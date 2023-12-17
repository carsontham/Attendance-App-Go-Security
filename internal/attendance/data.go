package attendance

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
)

// this function will write existing data to the data.json file
// data is saved by using json.Marshal and written to the data.json file
func SaveData(data any, filename string) error {
	fmt.Println(filename)
	filename = "internal/input/" + filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//writer := bufio.NewWriter(file)
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = file.Write(dataJson)
	if err != nil {
		return err
	}
	fmt.Println("Data saved to data.json")
	// return writer.Flush()
	return nil
}

// this function is called during the init() function
// this loads the data from the data.json file and saves it to the userList variable
// The data is read using bufio.Scanner and json.Unmarshal
func LoadData() error {
	var data UserList

	file, err := os.Open("internal/input/data.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	dataJson := scanner.Text()

	err = json.Unmarshal([]byte(dataJson), &data)
	if err != nil {
		return err
	}
	userList = data
	return nil
}

// this functions allows the data to be exported to a XML format
// output file will be named output.xml
func saveDataInXML(data any, filename string) error {
	xmlData, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	// Write the XML string to a file
	err = os.WriteFile("internal/output/output.xml", xmlData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveDataInJSON(data any, filename string) error {
	// fmt.Println(filename)
	filename = "internal/output/" + filename
	fmt.Println(filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//writer := bufio.NewWriter(file)
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = file.Write(dataJson)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Data saved to output.json")
	return nil
}
