package main

import (
        "encoding/json"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "os/signal"
        "path/filepath"
        "syscall"
        "time"

        "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
        "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
        log.Println("============ application-golang starts ============")

        err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
        if err != nil {
                log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
        }

        wallet, err := gateway.NewFileSystemWallet("wallet")
        if err != nil {
                log.Fatalf("Failed to create wallet: %v", err)
        }

        if !wallet.Exists("appUser") {
                err = populateWallet(wallet)
                if err != nil {
                        log.Fatalf("Failed to populate wallet contents: %v", err)
                }
        }

        ccpPath := filepath.Join(
                "..",
                "..",
                "test-network",
                "organizations",
                "peerOrganizations",
                "org1.example.com",
                "connection-org1.yaml",
        )

        gw, err := gateway.Connect(
                gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
                gateway.WithIdentity(wallet, "appUser"),
        )
        if err != nil {
                log.Fatalf("Failed to connect to gateway: %v", err)
        }
        defer gw.Close()

        http.HandleFunc("/initLedger", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodGet {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")
                result, err := InitLedgerTransaction(contract)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.Write(result)
        })

        http.HandleFunc("/medicines", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodGet {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")
                result, err := GetAllMedicinesTransaction(contract)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.Write(result)
        })

        http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }
                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        http.Error(w, "Failed to read request body", http.StatusBadRequest)
                        return
                }

                var medicine GetMedicine
                err = json.Unmarshal(body, &medicine)
                if err != nil {
                        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")
                result, err := GetMedicineTransaction(contract, medicine.ID)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.Write(result)
        })

        http.HandleFunc("/journey", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }
                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        http.Error(w, "Failed to read request body", http.StatusBadRequest)
                        return
                }

                var medicine GetMedicine
                err = json.Unmarshal(body, &medicine)
                if err != nil {
                        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")
                result, err := MedicineJourneyTransaction(contract, medicine.ID)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.Write(result)
        })

        http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }
                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        http.Error(w, "Failed to read request body", http.StatusBadRequest)
                        return
                }

                var medicine GetMedicine
                err = json.Unmarshal(body, &medicine)
                if err != nil {
                        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")
                result, err := GetMedicineHistoryTransaction(contract, medicine.ID)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.Write(result)
        })

        http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")

                // Read the request body
                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        http.Error(w, "Failed to read request body", http.StatusBadRequest)
                        return
                }

                // Parse the JSON request body into a Medicine struct
                var medicine Medicine4
                err = json.Unmarshal(body, &medicine)
                if err != nil {
                        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
                        return
                }

                currentTime := time.Now()
                formattedTime := currentTime.Format("2006-01-02 15:04:02 -0700 MST")

                medicine.TimeStamp = formattedTime

                result, err := CreateMedicineTransaction(contract,
                        medicine.ID, medicine.Name, medicine.Manufacturer, medicine.ManufactureDate, medicine.ExpiryDate,
                        medicine.BrandName, medicine.Composition, medicine.SenderID, medicine.ReceiverID,
                        medicine.DRAPNo, medicine.DosageForm, medicine.TimeStamp, medicine.Batch_No, medicine.JourneyCompleted)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        log.Println("Error submitting CreateMedicineTransaction:", err)
                        return
                }
                log.Println("success ", result)

                // Convert the medicine struct to JSON
                jsonResponse, err := json.Marshal(medicine)
                if err != nil {
                        http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
                        return
                }

                // Set the appropriate Content-Type header
                w.Header().Set("Content-Type", "application/json")

                // Write the response body
                w.Write(jsonResponse)
        })

        http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                contract := getContract(gw, "mychannel", "basic")

                // Read the request body
                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        http.Error(w, "Failed to read request body", http.StatusBadRequest)
                        return
                }
                // Parse the JSON request body into a Medicine struct
                var medicine Medicine4
                err = json.Unmarshal(body, &medicine)
                if err != nil {
                        http.Error(w, "Failed to parse request body", http.StatusBadRequest)
                        return
                }
                currentTime := time.Now()
                formattedTime := currentTime.Format("2006-01-02 15:04:02 -0700 MST")

                medicine.TimeStamp = formattedTime

                result, err := UpdateMedicineTransaction(contract,
                        medicine.ID, medicine.Name, medicine.Manufacturer, medicine.ManufactureDate, medicine.ExpiryDate,
                        medicine.BrandName, medicine.Composition, medicine.SenderID, medicine.ReceiverID,
                        medicine.DRAPNo, medicine.DosageForm, medicine.TimeStamp, medicine.Batch_No, medicine.JourneyCompleted)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                log.Println(result)
                log.Println("success")

        })

        http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
                if r.Method != http.MethodPost {
                        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                        return
                }

                // Shut down the server gracefully
                log.Println("Shutting down server...")
                err := shutdownServer()
                if err != nil {
                        log.Fatalf("Failed to shut down server: %v", err)
                }

                w.Write([]byte("Server is shutting down..."))
        })

        log.Println("Starting server on port 8080...")
        go func() {
                log.Fatal(http.ListenAndServe(":8080", nil))
        }()

        // Wait for interrupt signal to gracefully shut down the server
        waitForInterrupt()

        log.Println("============ application-golang stopped ============")
}

type Medicine4 struct {
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
        TimeStamp        string `json:"TimeStamp"`
        Batch_No         string `json:"Batch_No"`
        JourneyCompleted string `json:"JourneyCompleted"`
}

type GetMedicine struct {
        ID string `json:"ID"`
}

func getContract(gw *gateway.Gateway, channel, contractName string) *gateway.Contract {
        network, err := gw.GetNetwork(channel)
        if err != nil {
                log.Fatalf("Failed to get network: %v", err)
        }

        contract := network.GetContract(contractName)
        return contract
}

func InitLedgerTransaction(contract *gateway.Contract) ([]byte, error) {
        log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
        return contract.SubmitTransaction("InitLedger")
}

func GetAllMedicinesTransaction(contract *gateway.Contract) ([]byte, error) {
        log.Println("--> Evaluate Transaction: GetAllMedicines, function returns all the current assets on the ledger")
        return contract.EvaluateTransaction("GetAllMedicines")
}

func GetMedicineTransaction(contract *gateway.Contract, id string) ([]byte, error) {
        log.Println("--> Evaluate Transaction: GetMedicine, function returns all the current assets on the ledger")
        return contract.EvaluateTransaction("ReadMedicine", id)
}

func MedicineJourneyTransaction(contract *gateway.Contract, id string) ([]byte, error) {
        log.Println("--> Evaluate Transaction: GetMedicine, function returns Journey Completed Status")
        return contract.EvaluateTransaction("MedicineJourney", id)
}

func GetMedicineHistoryTransaction(contract *gateway.Contract, id string) ([]byte, error) {
        log.Println("--> Evaluate Transaction: GetMedicine, function returns all the current assets on the ledger")
        return contract.EvaluateTransaction("GetMedicineHistory", id)
}

func CreateMedicineTransaction(contract *gateway.Contract, id, name, manufacturer, manufactureDate, expiryDate,
        brandName, composition, senderID, receiverID,
        drapNo, dosageForm, description, batch_No, journeyCompleted string) ([]byte, error) {
        log.Println("--> Submit Transaction: CreateMedicine, creates a new medicine with the given details")

        response, err := contract.SubmitTransaction("CreateMedicine", id, name, manufacturer,
                manufactureDate,
                expiryDate,
                brandName, composition, senderID, receiverID,
                drapNo, dosageForm, description, batch_No, journeyCompleted)

        return response, err

}

func UpdateMedicineTransaction(contract *gateway.Contract, id, name, manufacturer, manufactureDate, expiryDate,
        brandName, composition, senderID, receiverID,
        drapNo, dosageForm, description, batch_No, journeyCompleted string) ([]byte, error) {
        log.Println("--> Submit Transaction: UpdateMedicine")
        return contract.SubmitTransaction("UpdateMedicine", id, name, manufacturer,
                manufactureDate,
                expiryDate,
                brandName, composition, senderID, receiverID,
                drapNo, dosageForm, description, batch_No, journeyCompleted)
}

func populateWallet(wallet *gateway.Wallet) error {
        log.Println("============ Populating wallet ============")
        credPath := filepath.Join(
                "..",
                "..",
                "test-network",
                "organizations",
                "peerOrganizations",
                "org1.example.com",
                "users",
                "User1@org1.example.com",
                "msp",
        )

        certPath := filepath.Join(credPath, "signcerts", "cert.pem")
        cert, err := ioutil.ReadFile(filepath.Clean(certPath))
        if err != nil {
                return err
        }

        keyDir := filepath.Join(credPath, "keystore")
        files, err := ioutil.ReadDir(keyDir)
        if err != nil {
                return err
        }

        for _, file := range files {
                keyPath := filepath.Join(keyDir, file.Name())
                key, err := ioutil.ReadFile(filepath.Clean(keyPath))
                if err != nil {
                        return err
                }

                err = wallet.Put("appUser", gateway.NewX509Identity("Org1MSP", string(cert), string(key)))
                if err != nil {
                        return err
                }
        }

        log.Println("============ Wallet populated ============")
        return nil
}

func shutdownServer() error {
        return syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func waitForInterrupt() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt, syscall.SIGTERM)
        <-c
}
