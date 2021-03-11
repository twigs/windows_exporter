package collector

import (
	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	registerCollector("cluster_csv_fs", newcluster_csv_fsCollector) // TODO: Add any perflib dependencies here
}

// A cluster_csv_fsCollector is a Prometheus collector for WMI Win32_PerfRawData_CsvFsPerfProvider_ClusterCSVFileSystem metrics
type cluster_csv_fsCollector struct {

	//General information
	VolumeState *prometheus.Desc
	CreateFile  *prometheus.Desc
	FilesOpened *prometheus.Desc
	Flushes     *prometheus.Desc

	//Counters indication how often (and why) the volume was paused
	VolumePauseCountDisk    *prometheus.Desc
	VolumePauseCountNetwork *prometheus.Desc
	VolumePauseCountOther   *prometheus.Desc
	VolumePauseCountTotal   *prometheus.Desc

	//all IO
	Reads            *prometheus.Desc
	Writes           *prometheus.Desc
	MetadataIO       *prometheus.Desc
	ReadLatency      *prometheus.Desc
	WriteLatency     *prometheus.Desc
	ReadQueueLength  *prometheus.Desc
	WriteQueueLength *prometheus.Desc

	//local IO
	IOReads               *prometheus.Desc
	IOWrites              *prometheus.Desc
	IOReadBytes           *prometheus.Desc
	IOWriteBytes          *prometheus.Desc
	IOReadLatency         *prometheus.Desc
	IOWriteLatency        *prometheus.Desc
	IOReadAvgQueueLength  *prometheus.Desc
	IOWriteAvgQueueLength *prometheus.Desc

	//redirected IO
	RedirectedReads                *prometheus.Desc
	RedirectedWrites               *prometheus.Desc
	RedirectedReadBytes            *prometheus.Desc
	RedirectedWriteBytes           *prometheus.Desc
	RedirectedReadLatency          *prometheus.Desc
	RedirectedWriteLatency         *prometheus.Desc
	RedirectedReadsAvgQueueLength  *prometheus.Desc
	RedirectedWritesAvgQueueLength *prometheus.Desc
}

func newcluster_csv_fsCollector() (Collector, error) {
	const subsystem = "cluster_csv_fs"
	return &cluster_csv_fsCollector{

		//General information
		VolumeState: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "volume_state"),
			"State of the volume. Volume can be in one of the following states. 0 - Init state. In that state all files are invalidated and all IOs except volume level IOs are failing. 1 - Paused state. In this state volume will pause any new IO and down-level state is cleaned. 2 - Draining state. In this state volume will pause any new IO, but down-level files are still opened and some down-level IOs might be still in process. 3 - Set down level state. In this state volume will pause any new IO. The down-level state is already reapplied. 4 - Active state. In this state all IO are proceeding as normal.(VolumeState)",
			[]string{"volume"},
			nil,
		),
		CreateFile: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "create_file"),
			"The number of files that were created on the volume.",
			[]string{"volume"},
			nil,
		),
		FilesOpened: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "files_opened"),
			"Number of files opened on this volume including the volume opens.",
			[]string{"volume"},
			nil,
		),
		Flushes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "flushes"),
			"The number of fushes that were performed on the volume.",
			[]string{"volume"},
			nil,
		),

		//Counters indication how often (and why) the volume was paused
		VolumePauseCountDisk: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "volume_pause_count_disk"),
			"Number of times this volume was paused due to an error from a disk. (VolumePauseCountDisk)",
			[]string{"volume"},
			nil,
		),
		VolumePauseCountNetwork: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "volume_pause_count_network"),
			"Number of times this volume was paused due to an error from a network. (VolumePauseCountNetwork)",
			[]string{"volume"},
			nil,
		),
		VolumePauseCountOther: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "volume_pause_count_other"),
			"Number of times this volume was paused for the reasons other than Direct IO or Network. (VolumePauseCountOther)",
			[]string{"volume"},
			nil,
		),
		VolumePauseCountTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "volume_pause_count_total"),
			"Number of times this volume was paused. (VolumePauseCountTotal)",
			[]string{"volume"},
			nil,
		),

		//all IO
		Reads: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "reads"),
			"The number of reads that were performed on the volume. (Reads)",
			[]string{"volume"},
			nil,
		),
		Writes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "writes"),
			"The number of writes that were performed on the volume. (Writes)",
			[]string{"volume"},
			nil,
		),
		MetadataIO: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "metadata_io"),
			"The number of metadata operations that were performed against the volume. (MetadataIO)",
			[]string{"volume"},
			nil,
		),
		ReadLatency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "read_latency"),
			"The average latency between the time a read request arrived to the file system and when it was completed. (ReadLatency)",
			[]string{"volume"},
			nil,
		),
		WriteLatency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "write_latency"),
			"The average latency between the time a write request arrived to the file system and when it was completed. (WriteLatency)",
			[]string{"volume"},
			nil,
		),
		ReadQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "read_queue_length"),
			"The number of reads outstanding on this volume (ReadQueueLength)",
			[]string{"volume"},
			nil,
		),
		WriteQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "write_queue_length"),
			"The number of writes outstanding on this volume. (WriteQueueLength)",
			[]string{"volume"},
			nil,
		),

		//local IO
		IOReads: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_reads"),
			"The number of reads that were performed directly from the volume.",
			[]string{"volume"},
			nil,
		),
		IOWrites: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_writes"),
			"The number of writes that were performed directly from the volume. (IOWrites)",
			[]string{"volume"},
			nil,
		),
		IOReadBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_read_bytes"),
			"The IO Read Bytes performance counter shows the number of bytes read directly from the volume.",
			[]string{"volume"},
			nil,
		),
		IOWriteBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_write_bytes"),
			"The number of bytes written directly to the volume (IOWriteBytes)",
			[]string{"volume"},
			nil,
		),
		IOReadLatency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_read_latency"),
			"The average latency between the time a read request is sent to the disk using Direct IO and when its response is received.",
			[]string{"volume"},
			nil,
		),
		IOWriteLatency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_write_latency"),
			"The average latency between the time a write request is sent to the disk using Direct IO and when its response is received. (IOWriteLatency)",
			[]string{"volume"},
			nil,
		),
		IOReadAvgQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_read_avg_queue_length"),
			"The average number of reads requests that were performed directly on the disk.",
			[]string{"volume"},
			nil,
		),
		IOWriteAvgQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "io_write_avg_queue_length"),
			"The average number of write requests that were performed directly on the disk.",
			[]string{"volume"},
			nil,
		),

		//redirected IO
		RedirectedReads: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_reads"),
			"The number of reads that were redirected to the volume over the network. (RedirectedReads)",
			[]string{"volume"},
			nil,
		),
		RedirectedWrites: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_writes"),
			"The number of writes that were redirected to the volume over the network. (RedirectedWrites)",
			[]string{"volume"},
			nil,
		),
		RedirectedReadBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_read_bytes"),
			"The number of bytes read that were redirected to the volume over the network. (RedirectedReadBytes)",
			[]string{"volume"},
			nil,
		),
		RedirectedWriteBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_write_bytes"),
			"The number of bytes written that were redirected to the volume over the network. (RedirectedWriteBytes)",
			[]string{"volume"},
			nil,
		),
		RedirectedReadLatency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_read_latency"),
			"The average latency between the time a read request was redirected to the volume through the network and when its response is received. (RedirectedReadLatency)",
			[]string{"volume"},
			nil,
		),
		RedirectedWriteLatency: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_write_latency"),
			"The average latency between the time a write request was redirected to the volume through the network and when its response is received. (RedirectedWriteLatency)",
			[]string{"volume"},
			nil,
		),
		RedirectedReadsAvgQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_reads_avg_queue_length"),
			"(RedirectedReadsAvgQueueLength)",
			[]string{"volume"},
			nil,
		),
		RedirectedWritesAvgQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "redirected_writes_avg_queue_length"),
			"(RedirectedWritesAvgQueueLength)",
			[]string{"volume"},
			nil,
		),
	}, nil
}

// Win32_PerfRawData_CsvFsPerfProvider_ClusterCSVFileSystem docs:
// - <add link to documentation here>
type Win32_PerfRawData_CsvFsPerfProvider_ClusterCSVFileSystem struct {
	Name string

	CreateFile                     uint64
	CreateFilePersec               uint64
	FilesInvalidatedDuringResume   uint64
	FilesInvalidatedOther          uint64
	FilesOpened                    uint32
	Flushes                        uint64
	FlushesPersec                  uint64
	IOReadAvgQueueLength           uint64
	IOReadBytes                    uint64
	IOReadBytesPersec              uint64
	IOReadLatency                  uint32
	IOReadQueueLength              uint64
	IOReads                        uint64
	IOReadsPersec                  uint64
	IOSingleReads                  uint64
	IOSingleReadsPersec            uint64
	IOSingleWrites                 uint64
	IOSingleWritesPersec           uint64
	IOSplitReads                   uint64
	IOSplitReadsPersec             uint64
	IOSplitWrites                  uint64
	IOSplitWritesPersec            uint64
	IOWriteAvgQueueLength          uint64
	IOWriteBytes                   uint64
	IOWriteBytesPersec             uint64
	IOWriteLatency                 uint32
	IOWriteQueueLength             uint64
	IOWrites                       uint64
	IOWritesPersec                 uint64
	MetadataIO                     uint64
	MetadataIOPersec               uint64
	ReadLatency                    uint32
	ReadQueueLength                uint64
	Reads                          uint64
	ReadsPersec                    uint64
	RedirectedReadBytes            uint64
	RedirectedReadBytesPersec      uint64
	RedirectedReadLatency          uint32
	RedirectedReadQueueLength      uint64
	RedirectedReads                uint64
	RedirectedReadsAvgQueueLength  uint64
	RedirectedReadsPersec          uint64
	RedirectedWriteBytes           uint64
	RedirectedWriteBytesPersec     uint64
	RedirectedWriteLatency         uint32
	RedirectedWriteQueueLength     uint64
	RedirectedWrites               uint64
	RedirectedWritesAvgQueueLength uint64
	RedirectedWritesPersec         uint64
	VolumePauseCountDisk           uint64
	VolumePauseCountNetwork        uint64
	VolumePauseCountOther          uint64
	VolumePauseCountTotal          uint64
	VolumeState                    uint32
	WriteLatency                   uint32
	WriteQueueLength               uint64
	Writes                         uint64
	WritesPersec                   uint64
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *cluster_csv_fsCollector) Collect(ctx *ScrapeContext, ch chan<- prometheus.Metric) error {
	var dst []Win32_PerfRawData_CsvFsPerfProvider_ClusterCSVFileSystem
	q := queryAll(&dst)
	if err := wmi.Query(q, &dst); err != nil {
		return err
	}

	for _, csvfs := range dst {

		if csvfs.Name == "_Total" {
			continue
		}

		//General information
		ch <- prometheus.MustNewConstMetric(c.VolumeState, prometheus.GaugeValue, float64(csvfs.VolumeState), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.CreateFile, prometheus.CounterValue, float64(csvfs.CreateFile), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.FilesOpened, prometheus.GaugeValue, float64(csvfs.FilesOpened), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.Flushes, prometheus.CounterValue, float64(csvfs.Flushes), csvfs.Name)

		//Counters indication how often (and why) the volume was paused
		ch <- prometheus.MustNewConstMetric(c.VolumePauseCountDisk, prometheus.CounterValue, float64(csvfs.VolumePauseCountDisk), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.VolumePauseCountNetwork, prometheus.CounterValue, float64(csvfs.VolumePauseCountNetwork), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.VolumePauseCountOther, prometheus.CounterValue, float64(csvfs.VolumePauseCountOther), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.VolumePauseCountTotal, prometheus.CounterValue, float64(csvfs.VolumePauseCountTotal), csvfs.Name)

		//all IO
		ch <- prometheus.MustNewConstMetric(c.Reads, prometheus.CounterValue, float64(csvfs.Reads), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.Writes, prometheus.CounterValue, float64(csvfs.Writes), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.MetadataIO, prometheus.CounterValue, float64(csvfs.MetadataIO), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.ReadLatency, prometheus.CounterValue, float64(csvfs.ReadLatency), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.WriteLatency, prometheus.CounterValue, float64(csvfs.WriteLatency), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.ReadQueueLength, prometheus.GaugeValue, float64(csvfs.ReadQueueLength), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.WriteQueueLength, prometheus.GaugeValue, float64(csvfs.WriteQueueLength), csvfs.Name)

		//local IO
		ch <- prometheus.MustNewConstMetric(c.IOReads, prometheus.CounterValue, float64(csvfs.IOReads), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOWrites, prometheus.CounterValue, float64(csvfs.IOWrites), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOReadBytes, prometheus.CounterValue, float64(csvfs.IOReadBytes), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOWriteBytes, prometheus.CounterValue, float64(csvfs.IOWriteBytes), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOReadLatency, prometheus.CounterValue, float64(csvfs.IOReadLatency), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOWriteLatency, prometheus.CounterValue, float64(csvfs.IOWriteLatency), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOReadAvgQueueLength, prometheus.CounterValue, float64(csvfs.IOReadAvgQueueLength), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.IOWriteAvgQueueLength, prometheus.CounterValue, float64(csvfs.IOWriteAvgQueueLength), csvfs.Name)

		//redirected IO
		ch <- prometheus.MustNewConstMetric(c.RedirectedReads, prometheus.CounterValue, float64(csvfs.RedirectedReads), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedWrites, prometheus.CounterValue, float64(csvfs.RedirectedWrites), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedReadBytes, prometheus.CounterValue, float64(csvfs.RedirectedReadBytes), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedWriteBytes, prometheus.CounterValue, float64(csvfs.RedirectedWriteBytes), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedReadLatency, prometheus.CounterValue, float64(csvfs.RedirectedReadLatency), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedWriteLatency, prometheus.CounterValue, float64(csvfs.RedirectedWriteLatency), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedReadsAvgQueueLength, prometheus.CounterValue, float64(csvfs.RedirectedReadsAvgQueueLength), csvfs.Name)
		ch <- prometheus.MustNewConstMetric(c.RedirectedWritesAvgQueueLength, prometheus.CounterValue, float64(csvfs.RedirectedWritesAvgQueueLength), csvfs.Name)
	}

	return nil
}
