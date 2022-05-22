package main

type BookMetadata struct {
	Title string `json:"title"`
	Level int    `json:"level"`
}

type Index struct {
	Keys Key `json:"keys"`
}

type Key struct {
	Kindle float32 `json:"kindle"`
	Audio  string  `json:"audio"`
}

func (k Index) LessThan(k2 Index) bool {
	kindleLessThan := false
	if k.Keys.Kindle < k2.Keys.Kindle {
		kindleLessThan = true
	}
	audioLessThan := false
	if k.Keys.Audio != "" && k2.Keys.Audio != "" {
		if k.Keys.Audio < k2.Keys.Audio {
			audioLessThan = true
		}
	}
	return kindleLessThan || audioLessThan
}

func (k Index) GreaterThan(k2 Index) bool {
	kindleLessThan := false
	if k.Keys.Kindle > k2.Keys.Kindle {
		kindleLessThan = true
	}
	audioLessThan := false
	if k.Keys.Audio != "" && k2.Keys.Audio != "" {
		if k.Keys.Audio > k2.Keys.Audio {
			audioLessThan = true
		}
	}
	return kindleLessThan || audioLessThan
}
