package main

type ValueMessage struct {
	Value  float32 `json:"value"`
	Client string  `json:"client"`
}
