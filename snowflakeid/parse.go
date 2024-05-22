package snowflakeid

func Parse(layout Layout, id int64) (epoch, datacenter, worker, sequence int64, err error) {
	_, _, datacenterMax, workerMax, sequenceMax, idShift, datacenterShift, ex := layout.Validate()
	if ex != nil {
		err = ex
		return
	}

	sequenceBits := int64(layout.SequenceBits)

	epoch = id>>idShift + int64(layout.Epoch)
	datacenter = id >> datacenterShift & datacenterMax
	worker = id >> sequenceBits & workerMax
	sequence = id & sequenceMax

	return
}
