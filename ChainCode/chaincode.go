package main

import (
        "encoding/json"
        "fmt"
        "log"

        "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing medicine data
type SmartContract struct {
        contractapi.Contract
}

// Medicine describes the details of a medicine
type Medicine struct {
        ID               string `json:"ID"`
        Name             string `json:"Name"`
        Manufacturer     string `json:"Manufacturer"`
        ManufactureDate  string `json:"ManufactureDate"`
        ExpiryDate       string `json:"ExpiryDate"`
        BrandName        string `json:"BrandName"`
        Composition      string `json:"Composition"`
        SenderID         string `json:"SenderId"`
        ReceiverID       string `json:"ReceiverId"`
        DRAPNo           string `json:"DrapNo"`
        DosageForm       string `json:"DosageForm"`
        TimeStamp      string `json:"TimeStamp"`
        Batch_No         string `json:"Batch_No"`
        JourneyCompleted string `json:"JourneyCompleted"`
}

// InitLedger adds a base set of medicines to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
        medicines := []Medicine{
                {
                        ID:               "1",
                        Name:             "Aspirin",
                        Manufacturer:     "ABC Pharmaceuticals",
                        ManufactureDate:  "2022-01-01",
                        ExpiryDate:       "2024-01-01",
                        BrandName:        "PharmaCorp",
                        Composition:      "Acetylsalicylic Acid",
                        SenderID:         "Sender1",
                        ReceiverID:       "Receiver1",
                        DRAPNo:           "1",
                        DosageForm:       "Tablet",
                        TimeStamp:        "Pain Reliever",
                        Batch_No:         "BatchNo1",
                        JourneyCompleted: "false",
                },
                {
                        ID:               "2",
                        Name:             "Paracetamol",
                        Manufacturer:     "XYZ Pharmaceuticals",
                        ManufactureDate:  "2022-02-01",
                        ExpiryDate:       "2024-02-01",
                        BrandName:        "HealthCare",
                        Composition:      "Paracetamol",
                        SenderID:         "Sender2",
                        ReceiverID:       "Receiver2",
                        DRAPNo:           "2",
                        DosageForm:       "Tablet",
                        TimeStamp:        "Fever Reducer",
                        Batch_No:         "BatchNo2",
                        JourneyCompleted: "false",
                },

                {
                        ID:               "3",
                        Name:             "Ibuprofen",
                        Manufacturer:     "PQR Pharmaceuticals",
                        ManufactureDate:  "2022-03-01",
                        ExpiryDate:       "2024-03-01",
                        BrandName:        "MediLife",
                        Composition:      "Ibuprofen",
                        SenderID:         "Sender3",
                        ReceiverID:       "Receiver3",
                        DRAPNo:           "3",
                        DosageForm:       "Capsule",
                        TimeStamp:        "Anti-inflammatory",
                        Batch_No:         "BatchNo3",
                        JourneyCompleted: "false",
                },

                {
                        ID:               "4",
                        Name:             "Amoxicillin",
                        Manufacturer:     "LMN Pharmaceuticals",
                        ManufactureDate:  "2022-04-01",
                        ExpiryDate:       "2024-04-01",
                        BrandName:        "PharmaMed",
                        Composition:      "Amoxicillin",
                        SenderID:         "Sender4",
                        ReceiverID:       "Receiver4",
                        DRAPNo:           "4",
                        DosageForm:       "Tablet",
                        TimeStamp:        "Antibiotic",
                        Batch_No:         "BatchNo4",
                        JourneyCompleted: "false",
                },
                {
                        ID:               "5",
                        Name:             "Omeprazole",
                        Manufacturer:     "EFG Pharmaceuticals",
                        ManufactureDate:  "2022-05-01",
                        ExpiryDate:       "2024-05-01",
                        BrandName:        "PharmaCare",
                        Composition:      "Omeprazole",
                        SenderID:         "Sender5",
                        ReceiverID:       "Receiver5",
                        DRAPNo:           "5",
                        DosageForm:       "Capsule",
                        TimeStamp:        "Acid Reducer",
                        Batch_No:         "BatchNo5",
                        JourneyCompleted: "false",
                },
                // Add more medicines here...
        }

        for _, medicine := range medicines {
                medicineJSON, err := json.Marshal(medicine)
                if err != nil {
                        return fmt.Errorf("failed to marshal medicine JSON: %v", err)
                }

                err = ctx.GetStub().PutState(medicine.ID, medicineJSON)
                if err != nil {
                        return fmt.Errorf("failed to put medicine in world state: %v", err)
                }
        }

        return nil
}

// CreateMedicine issues a new medicine to the world state with the given details.
func (s *SmartContract) CreateMedicine(ctx contractapi.TransactionContextInterface, id string,
        name string, manufacturer string, manufactureDate string,
        expiryDate string, brandName string, composition string, senderID string,
        receiverID string, drApNo string, dosageForm string, timeStamp string, batch_No string, journeyCompleted string) error {

        exists, err := s.MedicineExists(ctx, id)
        if err != nil {
                return fmt.Errorf("failed to check medicine existence: %v", err)
        }
        if exists {
                return fmt.Errorf("the medicine %s already exists", id)
        }

        medicine := Medicine{
                ID:               id,
                Name:             name,
                Manufacturer:     manufacturer,
                ManufactureDate:  manufactureDate,
                ExpiryDate:       expiryDate,
                BrandName:        brandName,
                Composition:      composition,
                SenderID:         senderID,
                ReceiverID:       receiverID,
                DRAPNo:           drApNo,
                DosageForm:       dosageForm,
                TimeStamp:        timeStamp,
                Batch_No:         batch_No,
                JourneyCompleted: journeyCompleted,
        }
        medicineJSON, err := json.Marshal(medicine)
        if err != nil {
                return fmt.Errorf("failed to marshal medicine JSON: %v", err)
        }

        err = ctx.GetStub().PutState(id, medicineJSON)
        if err != nil {
                return fmt.Errorf("failed to put medicine in world state: %v", err)
        }

        return nil
}

// ReadMedicine returns the medicine stored in the world state with the given id.
func (s *SmartContract) ReadMedicine(ctx contractapi.TransactionContextInterface, id string) (*Medicine, error) {
        medicineJSON, err := ctx.GetStub().GetState(id)
        if err != nil {
                return nil, fmt.Errorf("failed to read medicine from world state: %v", err)
        }
        if medicineJSON == nil {
                return nil, fmt.Errorf("the medicine %s does not exist", id)
        }

        var medicine Medicine
        err = json.Unmarshal(medicineJSON, &medicine)
        if err != nil {
                return nil, fmt.Errorf("failed to unmarshal medicine JSON: %v", err)
        }

        return &medicine, nil
}

// UpdateMedicine updates an existing medicine in the world state with the provided parameters.
func (s *SmartContract) UpdateMedicine(ctx contractapi.TransactionContextInterface,
        id string, name string, manufacturer string, manufactureDate string,
        expiryDate string, brandName string, composition string, senderID string,
        receiverID string, drApNo string, dosageForm string, timeStamp string, batch_No string, journeyCompleted string) error {
        exists, err := s.MedicineExists(ctx, id)
        if err != nil {
                return fmt.Errorf("failed to check medicine existence: %v", err)
        }
        if !exists {
                return fmt.Errorf("the medicine %s does not exist", id)
        }

        medicine := Medicine{
                ID:               id,
                Name:             name,
                Manufacturer:     manufacturer,
                ManufactureDate:  manufactureDate,
                ExpiryDate:       expiryDate,
                BrandName:        brandName,
                Composition:      composition,
                SenderID:         senderID,
                ReceiverID:       receiverID,
                DRAPNo:           drApNo,
                DosageForm:       dosageForm,
                TimeStamp:        timeStamp,
                Batch_No:         batch_No,
                JourneyCompleted: journeyCompleted,
        }
        medicineJSON, err := json.Marshal(medicine)
        if err != nil {
                return fmt.Errorf("failed to marshal medicine JSON: %v", err)
        }

        err = ctx.GetStub().PutState(id, medicineJSON)
        if err != nil {
                return fmt.Errorf("failed to put medicine in world state: %v", err)
        }

        return nil
}

// DeleteMedicine deletes a given medicine from the world state.
func (s *SmartContract) DeleteMedicine(ctx contractapi.TransactionContextInterface, id string) error {
        exists, err := s.MedicineExists(ctx, id)
        if err != nil {
                return fmt.Errorf("failed to check medicine existence: %v", err)
        }
        if !exists {
                return fmt.Errorf("the medicine %s does not exist", id)
        }

        err = ctx.GetStub().DelState(id)
        if err != nil {
                return fmt.Errorf("failed to delete medicine from world state: %v", err)
        }

        return nil
}

// MedicineExists returns true when a medicine with the given ID exists in the world state.
func (s *SmartContract) MedicineExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
        medicineJSON, err := ctx.GetStub().GetState(id)
        if err != nil {
                return false, fmt.Errorf("failed to read medicine from world state: %v", err)
        }

        return medicineJSON != nil, nil
}

// TransferMedicine updates the SenderId and RecieverId fields of a medicine with the given id in the world state, and returns the old owner.
func (s *SmartContract) TransferMedicine(ctx contractapi.TransactionContextInterface, id string,
        senderId string, receiverId string) (string, error) {
        medicine, err := s.ReadMedicine(ctx, id)

        if err != nil {
                return "", fmt.Errorf("failed to read medicine: %v", err)
        }

        oldSenderId := medicine.SenderID
        oldReceiverId := medicine.ReceiverID

        medicine.SenderID = senderId
        medicine.ReceiverID = receiverId

        medicineJSON, err := json.Marshal(medicine)
        if err != nil {
                return "", fmt.Errorf("failed to marshal medicine JSON: %v", err)
        }

        err = ctx.GetStub().PutState(id, medicineJSON)
        if err != nil {
                return "", fmt.Errorf("failed to put medicine in world state: %v", err)
        }

        return fmt.Sprintf("Previous SenderId: %s, Previous ReceiverId: %s", oldSenderId, oldReceiverId), nil
}

func (s *SmartContract) MedicineJourney(ctx contractapi.TransactionContextInterface, id string) (*Medicine, error) {
        medicine, err := s.ReadMedicine(ctx, id)

        if err != nil {
                return medicine, fmt.Errorf("failed to read medicine: %v", err)
        }

        medicine.JourneyCompleted = "true"

        medicineJSON, err := json.Marshal(medicine)
        if err != nil {
                return medicine, fmt.Errorf("failed to marshal medicine JSON: %v", err)
        }

        err = ctx.GetStub().PutState(id, medicineJSON)
        if err != nil {
                return medicine, fmt.Errorf("failed to put medicine in world state: %v", err)
        }

        return medicine, nil
}

// GetAllMedicines returns all medicines found in the world state.
func (s *SmartContract) GetAllMedicines(ctx contractapi.TransactionContextInterface) ([]*Medicine, error) {
        resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
        if err != nil {
                return nil, fmt.Errorf("failed to get medicines from world state: %v", err)
        }
        defer resultsIterator.Close()

        var medicines []*Medicine
        for resultsIterator.HasNext() {
                queryResponse, err := resultsIterator.Next()
                if err != nil {
                        return nil, fmt.Errorf("failed to iterate over medicines: %v", err)
                }

                var medicine Medicine
                err = json.Unmarshal(queryResponse.Value, &medicine)
                if err != nil {
                        return nil, fmt.Errorf("failed to unmarshal medicine JSON: %v", err)
                }
                medicines = append(medicines, &medicine)
        }

        return medicines, nil
}

// GetMedicineHistory returns the history of changes for a medicine with the given ID.
func (s *SmartContract) GetMedicineHistory(ctx contractapi.TransactionContextInterface, id string) ([]*Medicine, error) {
        resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
        if err != nil {
                return nil, fmt.Errorf("failed to get history for medicine: %v", err)
        }
        defer resultsIterator.Close()

        var history []*Medicine

        for resultsIterator.HasNext() {
                response, err := resultsIterator.Next()
                if err != nil {
                        return nil, fmt.Errorf("failed to iterate history for medicine: %v", err)
                }

                var medicine Medicine
                err = json.Unmarshal(response.Value, &medicine)
                if err != nil {
                        return nil, fmt.Errorf("failed to unmarshal history value for medicine: %v", err)
                }

                history = append(history, &medicine)
        }

        return history, nil
}

func main() {
        chaincode, err := contractapi.NewChaincode(&SmartContract{})
        if err != nil {
                log.Panicf("Error creating medicine-data chaincode: %v", err)
        }

        if err := chaincode.Start(); err != nil {
                log.Panicf("Error starting medicine-data chaincode: %v", err)
        }
}
