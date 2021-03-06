package service

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/watson-developer-cloud/go-sdk/personalityinsightsv3"
	"io/ioutil"
	"os"
)

type UserProfile struct {
	profile personalityinsightsv3.Profile
}

func (profile *UserProfile) Openness() int {

	value := profile.profile.Personality[0].Percentile
	intValue := int(*value * 100)
	return intValue
}

func (profile *UserProfile) Conscientiousness() int {

	value := *profile.profile.Personality[1].Percentile
	intValue := int(value * 100)
	return intValue
}
func (profile *UserProfile) Extraversion() int {

	value := profile.profile.Personality[2].Percentile
	intValue := int(*value * 100)
	return intValue
}

func (profile *UserProfile) Agreeableness() int {

	value := profile.profile.Personality[3].Percentile
	intValue := int(*value * 100)
	return intValue
}

func (profile *UserProfile) Neuroticism() int {

	value := profile.profile.Personality[4].Percentile
	intValue := int(*value * 100)
	return intValue
}

func (profile *UserProfile) SaveProfileAsJson(path string) error {
	fo, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create UserProfile save path")
	}

	defer fo.Close()
	encoder := json.NewEncoder(fo)

	err = encoder.Encode(profile.profile)
	if err != nil {
		return errors.Wrapf(err, "failed to create encode")
	}
	return nil
}

func (profile *UserProfile) LoadJsonAsProfile(path string) error {
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", path)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errors.Wrapf(err, "failed to load %s into Json", path)
	}
	json.Unmarshal(byteValue, &profile.profile)

	return nil
}
