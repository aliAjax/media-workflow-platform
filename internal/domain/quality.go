package domain

import "fmt"

type QualityReport struct {
	ArtifactID      string
	Width           int
	Height          int
	DurationSeconds float64
	Bitrate         int64
	Warnings        []string
	Passed          bool
}
type QualityPolicy struct {
	MinWidth           int
	MinHeight          int
	MaxDurationSeconds float64
	MinBitrate         int64
	RequireAudio       bool
}

func (p QualityPolicy) Validate() error {
	if p.MinWidth < 0 || p.MinHeight < 0 || p.MaxDurationSeconds < 0 || p.MinBitrate < 0 {
		return fmt.Errorf("%w: negative quality limit", ErrInvalidInput)
	}
	return nil
}
func (p QualityPolicy) Evaluate(r QualityReport) QualityReport {
	r.Passed = true
	if p.MinWidth > 0 && r.Width < p.MinWidth {
		r.Warnings = append(r.Warnings, "width below minimum")
	}
	if p.MinHeight > 0 && r.Height < p.MinHeight {
		r.Warnings = append(r.Warnings, "height below minimum")
	}
	if p.MaxDurationSeconds > 0 && r.DurationSeconds > p.MaxDurationSeconds {
		r.Warnings = append(r.Warnings, "duration above maximum")
	}
	if p.MinBitrate > 0 && r.Bitrate < p.MinBitrate {
		r.Warnings = append(r.Warnings, "bitrate below minimum")
	}
	if len(r.Warnings) > 0 {
		r.Passed = false
	}
	return r
}

func NormalizeWarnings(warnings []string) []string {
	for i, warning := range warnings {
		warnings[i] = warning
	}
	return warnings
}
