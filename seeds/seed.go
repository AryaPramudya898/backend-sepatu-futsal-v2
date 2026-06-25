package main

import (
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/config"
	"github.com/AryaPramudya898/backend-sepatu-futsal.git/models"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	godotenv.Load()
	config.InitDatabase()
	products := []models.Product{
		{
			Name:        "Ortuseight Catalyst Legion",
			Price:       399000,
			Category:    "Turf",
			Stock:       25,
			Description: "Sepatu futsal turf dengan grip kuat untuk lapangan sintetis",
			ImageURL:    "https://i.ibb.co.com/TBCqv3gy/1690005513.jpg",
		},
		{
			Name:        "Specs Accelerator Lightspeed",
			Price:       425000,
			Category:    "Turf",
			Stock:       20,
			Description: "Sepatu turf ringan untuk pergerakan cepat",
			ImageURL:    "https://kingofdribble.co.id/cdn/shop/files/1020298.jpg?v=1758617334",
		},
		{
			Name:        "Nike Mercurial Vapor TF",
			Price:       899000,
			Category:    "Turf",
			Stock:       12,
			Description: "Sepatu turf premium dengan kontrol bola maksimal",
			ImageURL:    "https://i.ibb.co.com/V0Nn67Tj/1-39.jpg",
		},
		{
			Name:        "Adidas Predator Edge TF",
			Price:       950000,
			Category:    "Turf",
			Stock:       10,
			Description: "Sepatu turf dengan stabilitas tinggi saat bermain",
			ImageURL:    "https://thumblr.uniid.it/product/238905/fceab34314bb.jpg?width=3840&format=webp&q=75",
		},
		{
			Name:        "Mizuno Morelia Neo TF",
			Price:       780000,
			Category:    "Turf",
			Stock:       14,
			Description: "Sepatu turf nyaman dengan bahan fleksibel",
			ImageURL:    "https://tha.mizuno.com/cdn/shop/files/SH_Q1GB249045_00_4b1d9609-3654-42ae-bc66-ffefb8184c5a.png?v=1730712855&width=1080",
		},
		{
			Name:        "Ortuseight Jogosala Indoor",
			Price:       389000,
			Category:    "Indoor",
			Stock:       22,
			Description: "Sepatu indoor dengan sol datar anti slip",
			ImageURL:    "https://i.ibb.co.com/m598Y8Ys/ortuseight-jogosala-lineage-whitegreen-1-yitwogcm.webp",
		},
		{
			Name:        "Specs Metasala Nativ",
			Price:       410000,
			Category:    "Indoor",
			Stock:       18,
			Description: "Sepatu indoor untuk kontrol dan passing akurat",
			ImageURL:    "https://i.ibb.co.com/FkP9VBpT/1647585409.jpg",
		},
		{
			Name:        "Nike Lunar Gato Indoor",
			Price:       1099000,
			Category:    "Indoor",
			Stock:       8,
			Description: "Sepatu indoor premium dengan bantalan empuk",
			ImageURL:    "https://static.nike.com/a/images/t_web_pdp_936_v2/f_auto,u_9ddf04c7-2a9a-4d76-add1-d15af8f0263d,c_scale,fl_relative,w_1.0,h_1.0,fl_layer_apply/ea8bd08c-4b0b-4861-896f-d8f4a7afcb49/NIKE+LUNARGATO+II.png",
		},
		{
			Name:        "Adidas Top Sala Indoor",
			Price:       699000,
			Category:    "Indoor",
			Stock:       11,
			Description: "Sepatu indoor dengan grip kuat di lantai halus",
			ImageURL:    "https://assets.adidas.com/images/w_600,f_auto,q_auto/c3aa9baa788e4b7c9817ea081e0fc4c8_9366/TOP_SALA_COMPETITION_II_Indoor_Football_Shoes_White_JP6980_01_00_standard.jpg",
		},
		{
			Name:        "Puma Future Z Indoor",
			Price:       759000,
			Category:    "Indoor",
			Stock:       9,
			Description: "Sepatu indoor fleksibel untuk manuver cepat",
			ImageURL:    "https://img.lazcdn.com/g/p/3698959fdf0528da692c86b9d89bb051.jpg_720x720q80.jpg",
		},
		{
			Name:        "Kelme Precision Elite",
			Price:       520000,
			Category:    "Turf",
			Stock:       16,
			Description: "Sepatu turf dengan kenyamanan tinggi dan grip stabil",
			ImageURL:    "https://kelme.com/cdn/shop/files/12869.jpg?v=1748010390&width=2048",
		},
		{
			Name:        "Joma Top Flex Turf",
			Price:       610000,
			Category:    "Turf",
			Stock:       13,
			Description: "Sepatu turf fleksibel untuk kontrol bola lebih baik",
			ImageURL:    "https://www.joma-sport.com/dw/image/v2/BFRV_PRD/on/demandware.static/-/Sites-joma-masterCatalog/default/dw01d3ec64/images/medium/TORS2601TF_1.jpg?sw=900&sh=900&sm=fit",
		},
		{
			Name:        "Munich Gresca Indoor",
			Price:       845000,
			Category:    "Indoor",
			Stock:       7,
			Description: "Sepatu indoor premium dengan daya cengkeram maksimal",
			ImageURL:    "https://www.munichsports.com/cdnassets/V6/SPORTS/3003003-01-L.jpg",
		},
		{
			Name:        "Asics Calcetto Indoor",
			Price:       575000,
			Category:    "Indoor",
			Stock:       15,
			Description: "Sepatu indoor ringan untuk permainan cepat",
			ImageURL:    "https://topsystem.id/api/product/600/1645586127.jpg",
		},
		{
			Name:        "Umbro Chaleira Pro",
			Price:       465000,
			Category:    "Indoor",
			Stock:       17,
			Description: "Sepatu indoor dengan bantalan nyaman untuk permainan lama",
			ImageURL:    "https://i.ibb.co.com/F4rbwYWH/ECFQtiy-Xs-AA1ym8.jpg",
		},
	}
	for _, p := range products {
		config.DB.Create(&p)
	}
	log.Printf("Seed berhasil: %d produk ditambahkan", len(products))
}
