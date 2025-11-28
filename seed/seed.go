package main

import (
	"context"
	"git-uts/config"
	"git-uts/model"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var places = []model.Location{
    // Jakarta
    {
        Name:        "Monas (Monumen Nasional)",
        Latitude:    -6.175387,
        Longitude:   106.827153,
        Description: "Landmark ikonik di pusat Jakarta.",
    },
    {
        Name:        "Istana Negara",
        Latitude:    -6.179394,
        Longitude:   106.827061,
        Description: "Istana resmi Presiden Republik Indonesia.",
    },
    {
        Name:        "Taman Mini Indonesia Indah",
        Latitude:    -6.316667,
        Longitude:   106.875000,
        Description: "Taman wisata budaya dan rekreasi.",
    },
    {
        Name:        "Kota Tua Jakarta",
        Latitude:    -6.132000,
        Longitude:   106.808333,
        Description: "Kawasan wisata sejarah di Jakarta Barat.",
    },
    {
        Name:        "Grand Indonesia Mall",
        Latitude:    -6.195411,
        Longitude:   106.820071,
        Description: "Salah satu pusat perbelanjaan mewah di Jakarta Pusat.",
    },
    // Bandung
    {
        Name:        "Gedung Sate",
        Latitude:    -6.917464,
        Longitude:   107.619110,
        Description: "Landmark dan kantor pemerintahan di Bandung.",
    },
    {
        Name:        "Kota Bandung",
        Latitude:    -6.917464, // Pusat Kota Bandung
        Longitude:   107.619110,
        Description: "Pusat pemerintahan dan kota utama di Jawa Barat.",
    },
    {
        Name:        "Trans Studio Bandung",
        Latitude:    -6.883333,
        Longitude:   107.583333,
        Description: "Taman hiburan indoor terbesar di Asia Tenggara.",
    },
    {
        Name:        "Dago Pakar",
        Latitude:    -6.850000,
        Longitude:   107.600000,
        Description: "Tempat wisata dan kuliner di kawasan Dago.",
    },
    {
        Name:        "Farm House Susu Lembang",
        Latitude:    -6.791667,
        Longitude:   107.616667,
        Description: "Tempat wisata dengan nuansa Eropa di Lembang.",
    },
}

func seedData(collection *mongo.Collection) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Hapus semua data yang ada di koleksi (opsional, uncomment jika ingin reset)
    // _, err := collection.DeleteMany(ctx, bson.M{})
    // if err != nil {
    //     return err
    // }
    // log.Println("Data lama dihapus.")

    // Masukkan data seed
    var docs []interface{}
    for i := range places {
        places[i].ID = primitive.NewObjectID() // Generate ID baru
        docs = append(docs, places[i])
    }

    result, err := collection.InsertMany(ctx, docs)
    if err != nil {
        return err
    }

    log.Printf("Berhasil memasukkan %d dokumen ke koleksi.\n", len(result.InsertedIDs))
    return nil
}

func main() {
    // Load environment variables dari .env
    if err := godotenv.Load("../.env"); err != nil {
        log.Println("⚠️  File .env tidak ditemukan, menggunakan environment variable yang ada")
    }

    log.Println("Menghubungkan ke MongoDB...")
    config.MongoConnect()

    db := config.DB // config.DB sudah bertipe *mongo.Database
    collection := db.Collection("location") // Ganti dengan nama koleksi kamu jika berbeda

    log.Println("Memulai seeding data...")
    err := seedData(collection)
    if err != nil {
        log.Fatal("Gagal seeding data: ", err)
    }

    log.Println("✅ Seeding data selesai.")
    os.Exit(0)
}