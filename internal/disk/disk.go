// Disk package contains abstract data-types to define disk-related entities.
//
// PartitionTable, Partition and Filesystem types are currently defined.
// All of them can be 1:1 converted to osbuild.QEMUAssemblerOptions.
package disk

import (
	"encoding/hex"
	"errors"
	"io"
	"math/rand"
	"sort"

	"github.com/google/uuid"
	osbuild "github.com/osbuild/osbuild-composer/internal/osbuild1"
	"github.com/osbuild/osbuild-composer/internal/osbuild2"
)

const (
	// Default sector size in bytes
	DefaultSectorSize = 512
)

type PartitionTable struct {
	Size       uint64 // Size of the disk (in bytes).
	UUID       string // Unique identifier of the partition table (GPT only).
	Type       string // Partition table type, e.g. dos, gpt.
	Partitions []Partition

	SectorSize uint64 // Sector size in bytes
}

type Partition struct {
	Start    uint64 // Start of the partition in sectors
	Size     uint64 // Size of the partition in sectors
	Type     string // Partition type, e.g. 0x83 for MBR or a UUID for gpt
	Bootable bool   // `Legacy BIOS bootable` (GPT) or `active` (DOS) flag
	// ID of the partition, dos doesn't use traditional UUIDs, therefore this
	// is just a string.
	UUID string
	// If nil, the partition is raw; It doesn't contain a filesystem.
	Filesystem *Filesystem
}

type Filesystem struct {
	Type string
	// ID of the filesystem, vfat doesn't use traditional UUIDs, therefore this
	// is just a string.
	UUID       string
	Label      string
	Mountpoint string
	// The fourth field of fstab(5); fs_mntops
	FSTabOptions string
	// The fifth field of fstab(5); fs_freq
	FSTabFreq uint64
	// The sixth field of fstab(5); fs_passno
	FSTabPassNo uint64
}

// Convert the given bytes to the number of sectors.
func (pt *PartitionTable) BytesToSectors(size uint64) uint64 {
	sectorSize := pt.SectorSize
	if sectorSize == 0 {
		sectorSize = DefaultSectorSize
	}
	return size / sectorSize
}

// Convert the given number of sectors to bytes.
func (pt *PartitionTable) SectorsToBytes(size uint64) uint64 {
	sectorSize := pt.SectorSize
	if sectorSize == 0 {
		sectorSize = DefaultSectorSize
	}
	return size * sectorSize
}

// Clone the partition table (deep copy).
func (pt *PartitionTable) Clone() *PartitionTable {
	if pt == nil {
		return nil
	}

	var partitions []Partition
	for _, p := range pt.Partitions {
		p.Filesystem = p.Filesystem.Clone()
		partitions = append(partitions, p)
	}
	return &PartitionTable{
		Size:       pt.Size,
		UUID:       pt.UUID,
		Type:       pt.Type,
		Partitions: partitions,

		SectorSize: pt.SectorSize,
	}
}

// Converts PartitionTable to osbuild.QEMUAssemblerOptions that encode
// the same partition table.
func (pt *PartitionTable) QEMUAssemblerOptions() osbuild.QEMUAssemblerOptions {
	var partitions []osbuild.QEMUPartition
	for _, p := range pt.Partitions {
		partitions = append(partitions, p.QEMUPartition())
	}

	return osbuild.QEMUAssemblerOptions{
		Size:       pt.Size,
		PTUUID:     pt.UUID,
		PTType:     pt.Type,
		Partitions: partitions,
	}
}

// Generates org.osbuild.fstab stage options from this partition table.
func (pt *PartitionTable) FSTabStageOptions() *osbuild.FSTabStageOptions {
	var options osbuild.FSTabStageOptions
	for _, p := range pt.Partitions {
		fs := p.Filesystem
		if fs == nil {
			continue
		}

		options.AddFilesystem(fs.UUID, fs.Type, fs.Mountpoint, fs.FSTabOptions, fs.FSTabFreq, fs.FSTabPassNo)
	}

	// sort the entries by PassNo to maintain backward compatibility
	sort.Slice(options.FileSystems, func(i, j int) bool {
		return options.FileSystems[i].PassNo < options.FileSystems[j].PassNo
	})

	return &options
}

// Generates org.osbuild.fstab stage options from this partition table.
func (pt *PartitionTable) FSTabStageOptionsV2() *osbuild2.FSTabStageOptions {
	var options osbuild2.FSTabStageOptions
	for _, p := range pt.Partitions {
		fs := p.Filesystem
		if fs == nil {
			continue
		}

		options.AddFilesystem(fs.UUID, fs.Type, fs.Mountpoint, fs.FSTabOptions, fs.FSTabFreq, fs.FSTabPassNo)
	}

	// sort the entries by PassNo to maintain backward compatibility
	sort.Slice(options.FileSystems, func(i, j int) bool {
		return options.FileSystems[i].PassNo < options.FileSystems[j].PassNo
	})

	return &options
}

// Returns the root partition (the partition whose filesystem has / as
// a mountpoint) of the partition table. Nil is returned if there's no such
// partition.
func (pt *PartitionTable) RootPartition() *Partition {
	for idx, p := range pt.Partitions {
		if p.Filesystem == nil {
			continue
		}

		if p.Filesystem.Mountpoint == "/" {
			return &pt.Partitions[idx]
		}
	}

	return nil
}

// Returns the /boot partition (the partition whose filesystem has /boot as
// a mountpoint) of the partition table. Nil is returned if there's no such
// partition.
func (pt *PartitionTable) BootPartition() *Partition {
	for _, p := range pt.Partitions {
		if p.Filesystem == nil {
			continue
		}

		if p.Filesystem.Mountpoint == "/boot" {
			return &p
		}
	}

	return nil
}

// Returns the index of the boot partition: the partition whose filesystem has
// /boot as a mountpoint.  If there is no explicit boot partition, the root
// partition is returned.
// If neither boot nor root partitions are found, returns -1.
func (pt *PartitionTable) BootPartitionIndex() int {
	// find partition with '/boot' mountpoint and fallback to '/'
	rootIdx := -1
	for idx, part := range pt.Partitions {
		if part.Filesystem == nil {
			continue
		}
		if part.Filesystem.Mountpoint == "/boot" {
			return idx
		} else if part.Filesystem.Mountpoint == "/" {
			rootIdx = idx
		}
	}
	return rootIdx
}

// StopIter is used as a return value from iterator function to indicate
// the iteration should not continue. Not an actual error and thus not
// returned by iterator function.
var StopIter = errors.New("stop the iteration")

// ForEachFileSystemFunc is a type of function called by ForEachFilesystem
// to iterate over every filesystem in the partition table.
//
// If the function returns an error, the iteration stops.
type ForEachFileSystemFunc func(fs *Filesystem) error

// Iterates over all filesystems in the partition table and calls the
// callback on each one. The iteration continues as long as the callback
// does not return an error.
func (pt *PartitionTable) ForEachFilesystem(cb ForEachFileSystemFunc) error {
	for _, part := range pt.Partitions {
		if part.Filesystem == nil {
			continue
		}

		if err := cb(part.Filesystem); err != nil {
			if err == StopIter {
				return nil
			}
			return err
		}
	}

	return nil
}

// Returns the Filesystem instance for a given mountpoint, if it exists.
func (pt *PartitionTable) FindFilesystemForMountpoint(mountpoint string) *Filesystem {
	var res *Filesystem
	_ = pt.ForEachFilesystem(func(fs *Filesystem) error {
		if fs.Mountpoint == mountpoint {
			res = fs
			return StopIter
		}

		return nil
	})

	return res
}

// Returns if the partition table contains a filesystem with the given
// mount point.
func (pt *PartitionTable) ContainsMountpoint(mountpoint string) bool {
	return pt.FindFilesystemForMountpoint(mountpoint) != nil
}

// Returns the Filesystem instance that corresponds to the root
// filesystem, i.e. the filesystem whose mountpoint is '/'.
func (pt *PartitionTable) RootFilesystem() *Filesystem {
	return pt.FindFilesystemForMountpoint("/")
}

// Returns the Filesystem instance that corresponds to the boot
// filesystem, i.e. the filesystem whose mountpoint is '/boot',
// if /boot is on a separate partition, otherwise nil
func (pt *PartitionTable) BootFilesystem() *Filesystem {
	return pt.FindFilesystemForMountpoint("/boot")
}

// Create a new filesystem within the partition table at the given mountpoint
// with the given minimum size in sectors.
func (pt *PartitionTable) CreateFilesystem(mountpoint string, size uint64) {
	filesystem := Filesystem{
		Type:         "xfs",
		Mountpoint:   mountpoint,
		FSTabOptions: "defaults",
		FSTabFreq:    0,
		FSTabPassNo:  0,
	}

	partition := Partition{
		Size:       size,
		Filesystem: &filesystem,
	}

	if pt.Type == "gpt" {
		partition.Type = FilesystemDataGUID
	}

	pt.Partitions = append(pt.Partitions, partition)
}

// Generate all needed UUIDs for all the partiton and filesystems
//
// Will not overwrite existing UUIDs and only generate UUIDs for
// partitions if the layout is GPT.
func (pt *PartitionTable) GenerateUUIDs(rng *rand.Rand) {
	_ = pt.ForEachFilesystem(func(fs *Filesystem) error {
		if fs.UUID == "" {
			fs.UUID = uuid.Must(newRandomUUIDFromReader(rng)).String()
		}
		return nil
	})

	// if this is a MBR partition table, there is no need to generate
	// uuids for the partitions themselves
	if pt.Type != "gpt" {
		return
	}

	for idx, part := range pt.Partitions {
		if part.UUID == "" {
			pt.Partitions[idx].UUID = uuid.Must(newRandomUUIDFromReader(rng)).String()
		}
	}
}

// dynamically calculate and update the start point
// for each of the existing partitions
// return the updated start point
func (pt *PartitionTable) updatePartitionStartPointOffsets(start uint64) uint64 {
	var rootIdx = -1
	for i := range pt.Partitions {
		partition := &pt.Partitions[i]
		if partition.Filesystem != nil && partition.Filesystem.Mountpoint == "/" {
			rootIdx = i
			continue
		}
		partition.Start = start
		start += partition.Size
	}
	pt.Partitions[rootIdx].Start = start
	return start
}

func (pt *PartitionTable) getPartitionTableSize() uint64 {
	var size uint64
	for _, p := range pt.Partitions {
		size += p.Size
	}
	return size
}

// Converts Partition to osbuild.QEMUPartition that encodes the same partition.
func (p *Partition) QEMUPartition() osbuild.QEMUPartition {
	var fs *osbuild.QEMUFilesystem
	if p.Filesystem != nil {
		f := p.Filesystem.QEMUFilesystem()
		fs = &f
	}
	return osbuild.QEMUPartition{
		Start:      p.Start,
		Size:       p.Size,
		Type:       p.Type,
		Bootable:   p.Bootable,
		UUID:       p.UUID,
		Filesystem: fs,
	}
}

// Filesystem related functions

// Clone the filesystem structure
func (fs *Filesystem) Clone() *Filesystem {
	if fs == nil {
		return nil
	}

	return &Filesystem{
		Type:         fs.Type,
		UUID:         fs.UUID,
		Label:        fs.Label,
		Mountpoint:   fs.Mountpoint,
		FSTabOptions: fs.FSTabOptions,
		FSTabFreq:    fs.FSTabFreq,
		FSTabPassNo:  fs.FSTabPassNo,
	}
}

// Converts Filesystem to osbuild.QEMUFilesystem that encodes the same fs.
func (fs *Filesystem) QEMUFilesystem() osbuild.QEMUFilesystem {
	return osbuild.QEMUFilesystem{
		Type:       fs.Type,
		UUID:       fs.UUID,
		Label:      fs.Label,
		Mountpoint: fs.Mountpoint,
	}
}

// uuid generator helpers

// GeneratesnewRandomUUIDFromReader generates a new random UUID (version
// 4 using) via the given random number generator.
func newRandomUUIDFromReader(r io.Reader) (uuid.UUID, error) {
	var id uuid.UUID
	_, err := io.ReadFull(r, id[:])
	if err != nil {
		return uuid.Nil, err
	}
	id[6] = (id[6] & 0x0f) | 0x40 // Version 4
	id[8] = (id[8] & 0x3f) | 0x80 // Variant is 10
	return id, nil
}

// NewRandomVolIDFromReader creates a random 32 bit hex string to use as a
// volume ID for FAT filesystems
func NewRandomVolIDFromReader(r io.Reader) (string, error) {
	volid := make([]byte, 4)
	_, err := r.Read(volid)
	return hex.EncodeToString(volid), err
}
