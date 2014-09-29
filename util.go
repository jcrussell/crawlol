package main

import "encoding/json"

// Hacky function to Marshal a MatchDetail into a MarshaledMatchDetail.
// Marshals all the fields that aren't basic types like strings and integers.
func (v *MatchDetail) Marshal() (*MarshaledMatchDetail, error) {
	var buf []byte
	var err error
	e := MarshaledMatchDetail{(*v).BaseMatchDetail, nil, nil, nil, nil}

	if buf, err = json.Marshal(v.ParticipantIdentities); err != nil {
		return nil, err
	}
	e.ParticipantIdentities = buf

	if buf, err = json.Marshal(v.Participants); err != nil {
		return nil, err
	}
	e.Participants = buf

	if buf, err = json.Marshal(v.Teams); err != nil {
		return nil, err
	}
	e.Teams = buf

	if buf, err = json.Marshal(v.Timeline); err != nil {
		return nil, err
	}
	e.Timeline = buf

	return &e, nil
}

// Hacky function to Unmarshal a MarshaledMatchDetail into a MatchDetail.
// Unmarshals all the fields that aren't basic types like strings and integers.
func (v *MarshaledMatchDetail) Unmarshal() (*MatchDetail, error) {
	e := MatchDetail{(*v).BaseMatchDetail, nil, nil, nil, Timeline{}}

	if err := json.Unmarshal(v.ParticipantIdentities, &e.ParticipantIdentities); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.Participants, &e.Participants); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.Teams, &e.Teams); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(v.Timeline, &e.Timeline); err != nil {
		return nil, err
	}

	return &e, nil
}
