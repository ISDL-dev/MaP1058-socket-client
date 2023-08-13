package model

type Status string

const (
	Acq         Status = "Acq"
	End         Status = "End"
	ReCalcEEG   Status = "ReCalcEEG"
	ReCalcHFLF  Status = "ReCalcHFLF"
	ReCalcEOG_V Status = "ReCalcEOG_V"
	ReCalcAll   Status = "ReCalcAll"
)
