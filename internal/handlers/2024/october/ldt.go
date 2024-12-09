package october

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var regionTranslations = map[string]string{
	"Новгородская область":                "Novgorod Region",
	"Владимирская область":                "Vladimir Region",
	"Мурманская область":                  "Murmansk Region",
	"Свердловская область":                "Sverdlovsk Region",
	"Калининградская область":             "Kaliningrad Region",
	"Тульская область":                    "Tula Region",
	"Рязанская область":                   "Ryazan Region",
	"Ярославская область":                 "Yaroslavl Region",
	"Воронежская область":                 "Voronezh Region",
	"Приморский край":                     "Primorsky Krai",
	"Чувашия Республика":                  "Chuvashia Republic",
	"Москва":                              "Moscow",
	"Сахалинская область":                 "Sakhalin Region",
	"Кировская область":                   "Kirov Region",
	"Белгородская область":                "Belgorod Region",
	"Красноярский край":                   "Krasnoyarsk Krai",
	"Новосибирская область":               "Novosibirsk Region",
	"Башкортостан Республика":             "Bashkortostan Republic",
	"Ненецкий автономный округ":           "Nenets Autonomous Okrug",
	"Чукотский автономный округ":          "Chukotka Autonomous Okrug",
	"Тамбовская область":                  "Tambov Region",
	"Чеченская Республика":                "Chechen Republic",
	"Коми Республика":                     "Komi Republic",
	"Алтайский край":                      "Altai Krai",
	"Татарстан Республика":                "Tatarstan Republic",
	"Иркутская область":                   "Irkutsk Region",
	"Северная Осетия Республика":          "North Ossetia Republic",
	"Ингушетия Республика":                "Ingushetia Republic",
	"Крым Республика":                     "Crimea Republic",
	"Магаданская область":                 "Magadan Region",
	"Саха (Якутия) Республика":            "Sakha (Yakutia) Republic",
	"Липецкая область":                    "Lipetsk Region",
	"Смоленская область":                  "Smolensk Region",
	"Орловская область":                   "Oryol Region",
	"Санкт-Петербург":                     "Saint Petersburg",
	"Луганская Народная Республика":       "Luhansk People's Republic",
	"Хакасия Республика":                  "Khakassia Republic",
	"Саратовская область":                 "Saratov Region",
	"Донецкая Народная Республика":        "Donetsk People's Republic",
	"Архангельская область":               "Arkhangelsk Region",
	"Нижегородская область":               "Nizhny Novgorod Region",
	"Волгоградская область":               "Volgograd Region",
	"Курская область":                     "Kursk Region",
	"Пензенская область":                  "Penza Region",
	"Тверская область":                    "Tver Region",
	"Челябинская область":                 "Chelyabinsk Region",
	"Московская область":                  "Moscow Region",
	"Забайкальский край":                  "Zabaykalsky Krai",
	"Ямало-Ненецкий автономный округ":     "Yamalo-Nenets Autonomous Okrug",
	"Брянская область":                    "Bryansk Region",
	"Курганская область":                  "Kurgan Region",
	"Удмуртия Республика":                 "Udmurtia Republic",
	"Самарская область":                   "Samara Region",
	"Калмыкия Республика":                 "Kalmykia Republic",
	"Ханты-Мансийский автономный округ":   "Khanty-Mansi Autonomous Okrug",
	"Адыгея Республика":                   "Adygea Republic",
	"Амурская область":                    "Amur Region",
	"Томская область":                     "Tomsk Region",
	"Тыва Республика":                     "Tuva Republic",
	"Кабардино-Балкария Республика":       "Kabardino-Balkaria Republic",
	"Астраханская область":                "Astrakhan Region",
	"Ивановская область":                  "Ivanovo Region",
	"Псковская область":                   "Pskov Region",
	"Карелия Республика":                  "Karelia Republic",
	"Севастополь":                         "Sevastopol",
	"Вологодская область":                 "Vologda Region",
	"Тюменская область":                   "Tyumen Region",
	"Оренбургская область":                "Orenburg Region",
	"Марий-Эл Республика":                 "Mari El Republic",
	"Ростовская область":                  "Rostov Region",
	"Краснодарский край":                  "Krasnodar Krai",
	"Алтай Республика":                    "Altai Republic",
	"Херсонская область":                  "Kherson Region",
	"Костромская область":                 "Kostroma Region",
	"Камчатский край":                     "Kamchatka Krai",
	"Омская область":                      "Omsk Region",
	"Запорожская область":                 "Zaporizhzhia Region",
	"Ленинградская область":               "Leningrad Region",
	"Ульяновская область":                 "Ulyanovsk Region",
	"Дагестан Республика":                 "Dagestan Republic",
	"Калужская область":                   "Kaluga Region",
	"Кемеровская область":                 "Kemerovo Region",
	"Пермский край":                       "Perm Krai",
	"Мордовия Республика":                 "Mordovia Republic",
	"Хабаровский край":                    "Khabarovsk Krai",
	"Еврейский автономный округ":          "Jewish Autonomous Okrug",
	"Карачаево-Черкессия Республика":      "Karachay-Cherkessia Republic",
	"Ставропольский край":                 "Stavropol Krai",
	"Бурятия Республика":                  "Buryatia Republic",
	"Центральный Федеральный Округ":       "Central",
	"Северо-Западный Федеральный Округ":   "North West",
	"Южный Федеральный Округ":             "South",
	"Северо-Кавказский Федеральный Округ": "North Caucasian",
	"Приволжский Федеральный Округ":       "Volga",
	"Уральский Федеральный Округ":         "Ural",
	"Сибирский Федеральный Округ":         "Siberia",
	"Дальневосточный Федеральный Округ":   "Far East",
}

// Мапа с переводами федеральных округов на английский
var districtTranslations = map[string]string{
	"Центральный Федеральный Округ":       "Central",
	"Северо-Западный Федеральный Округ":   "North West",
	"Южный Федеральный Округ":             "South",
	"Северо-Кавказский Федеральный Округ": "North Caucasian",
	"Приволжский Федеральный Округ":       "Volga",
	"Уральский Федеральный Округ":         "Ural",
	"Сибирский Федеральный Округ":         "Siberia",
	"Дальневосточный Федеральный Округ":   "Far East",
}

func TenMonth2024Ldt(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName  string `json:"region_name"`
		DONGFENG    int    `json:"dongfeng"`
		FOTON       int    `json:"foton"`
		GAZ         int    `json:"gaz"`
		ISUZU       int    `json:"isuzu"`
		JAC         int    `json:"jac"`
		KAMAZ       int    `json:"kamaz"`
		OTHER       int    `json:"other"`
		TotalMarket int    `json:"total"`
	}

	// Ответ с данными по анализу грузовиков
	type TruckAnalyticsResponse struct {
		Data  *orderedmap.OrderedMap[string, []*TruckAnalytics] `json:"data"`
		Error string                                            `json:"error,omitempty"`
	}

	// SQL запрос для получения данных
	query := `
WITH base_data AS (
    SELECT 
        "Federal_district",
        "Region",  
        "Brand",
        "City",
        SUM("Quantity") AS total_sales
    FROM ldt_3_5_12_truck_analytics_10_2024
    WHERE "Brand" IN ('DONGFENG', 'FOTON', 'GAZ', 'ISUZU', 'JAC', 'KAMAZ')
    GROUP BY "Federal_district", "Region", "City", "Brand"
),
federal_totals AS (
    SELECT 
        "Federal_district",
        "Region",
        'TOTAL' AS "City",
        "Brand",
        SUM(total_sales) AS total_sales
    FROM base_data
    GROUP BY "Federal_district", "Region", "Brand"
),
other_brands AS (
    SELECT 
        "Federal_district",
        "Region",  
        'OTHER' AS "Brand",
        'OTHER' AS "City",
        SUM("Quantity") AS total_sales
    FROM ldt_3_5_12_truck_analytics_10_2024
    WHERE "Brand" NOT IN ('DONGFENG', 'FOTON', 'GAZ', 'ISUZU', 'JAC', 'KAMAZ')
    GROUP BY "Federal_district", "Region"
),
combined_data AS (
    SELECT "Federal_district", "Region", "City", "Brand", total_sales
    FROM base_data
    UNION ALL
    SELECT "Federal_district", "Region", "City", "Brand", total_sales
    FROM federal_totals
    UNION ALL
    SELECT "Federal_district", "Region", "City", "Brand", total_sales
    FROM other_brands
),
pivoted_data AS (
    SELECT 
        "Region" AS Oblast_name,
        "Federal_district", 
        MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) AS DONGFENG,
        MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) AS FOTON,
        MAX(CASE WHEN "Brand" = 'GAZ' THEN total_sales END) AS GAZ,
        MAX(CASE WHEN "Brand" = 'ISUZU' THEN total_sales END) AS ISUZU,
        MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) AS JAC,
        MAX(CASE WHEN "Brand" = 'KAMAZ' THEN total_sales END) AS KAMAZ,
        MAX(CASE WHEN "Brand" = 'OTHER' THEN total_sales END) AS OTHER,
        COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'GAZ' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'ISUZU' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'KAMAZ' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'OTHER' THEN total_sales END), 0) AS TOTAL
    FROM combined_data
    GROUP BY "Region", "Federal_district"
),
federal_totals_by_region AS (
    SELECT 
        "Federal_district",
        "Federal_district" AS Oblast_name,  -- Use the Federal_district name in the total row
        SUM(DONGFENG) AS DONGFENG,
        SUM(FOTON) AS FOTON,
        SUM(GAZ) AS GAZ,
        SUM(ISUZU) AS ISUZU,
        SUM(JAC) AS JAC,
        SUM(KAMAZ) AS KAMAZ,
        SUM(OTHER) AS OTHER,
        SUM(TOTAL) AS TOTAL
    FROM pivoted_data
    GROUP BY "Federal_district"
),
final_data AS (
    SELECT Oblast_name, "Federal_district", DONGFENG, FOTON, GAZ, ISUZU, JAC, KAMAZ, OTHER, TOTAL
    FROM pivoted_data
    UNION ALL
    SELECT Oblast_name, "Federal_district", DONGFENG, FOTON, GAZ, ISUZU, JAC, KAMAZ, OTHER, TOTAL
    FROM federal_totals_by_region
)
SELECT 
    "Federal_district",  -- Federal_district comes first
    Oblast_name,  -- Oblast_name comes second
    COALESCE(DONGFENG, 0) AS DONGFENG,
    COALESCE(FOTON, 0) AS FOTON,
    COALESCE(GAZ, 0) AS GAZ,
    COALESCE(ISUZU, 0) AS ISUZU,
    COALESCE(JAC, 0) AS JAC,
    COALESCE(KAMAZ, 0) AS KAMAZ,
    COALESCE(OTHER, 0) AS OTHER,
    COALESCE(TOTAL, 0) AS TOTAL
FROM final_data
ORDER BY 
    "Federal_district",  -- Sort by Federal_district
    CASE 
        WHEN Oblast_name = "Federal_district" THEN 1  
        ELSE 0
    END,
    Oblast_name;
    `

	// Соединение с базой данных
	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Can't connect to database"})
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to execute query: " + err.Error()})
		return
	}
	defer rows.Close()

	// Ordered map to maintain custom order for districts
	dataByDistrict := orderedmap.New[string, []*TruckAnalytics]()

	// Define the custom order of districts
	customOrder := []string{
		"Central",
		"North West",
		"Volga",
		"South",
		"North Caucasian",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Prepopulate the map with the desired order
	for _, district := range customOrder {
		dataByDistrict.Set(district, []*TruckAnalytics{})
	}

	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.DONGFENG,
			&ta.FOTON,
			&ta.GAZ,
			&ta.ISUZU,
			&ta.JAC,
			&ta.KAMAZ,
			&ta.OTHER,
			&ta.TotalMarket,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to scan row: " + err.Error()})
			return
		}

		// Перевод федерального округа, если есть в маппинге
		if translatedDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = translatedDistrict
		}

		// Перевод региона, если есть в маппинге
		if translatedRegion, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = translatedRegion
		}

		// Рассчитываем общий рынок
		ta.TotalMarket = nullToZero(&ta.DONGFENG) + nullToZero(&ta.FOTON) + nullToZero(&ta.GAZ) + nullToZero(&ta.JAC) + nullToZero(&ta.ISUZU) + nullToZero(&ta.KAMAZ) + nullToZero(&ta.OTHER)

		// Add data to the appropriate district
		if existing, exists := dataByDistrict.Get(federalDistrict); exists {
			dataByDistrict.Set(federalDistrict, append(existing, &ta))
		}
	}

	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Error iterating over rows: " + err.Error()})
		return
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

type DistrictData struct {
	RegionName  string `json:"region_name"`
	DONGFENG    *int   `json:"dongfeng"`
	FOTON       *int   `json:"foton"`
	GAZ         *int   `json:"gaz"`
	ISUZU       *int   `json:"isuzu"`
	JAC         *int   `json:"jac"`
	KAMAZ       *int   `json:"kamaz"`
	OTHER       *int   `json:"other"`
	TotalMarket int    `json:"total"`
}

func TenMonth2024LDTTotal(ctx *gin.Context) {
	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictData] `json:"data"`
		Error string                                         `json:"error,omitempty"`
	}

	var districtTranslations = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far East",
		"Приволжский Федеральный Округ":       "Volga",
		"Северо-Западный Федеральный Округ":   "North West",
		"Северо-Кавказский Федеральный Округ": "North Caucasian",
		"Сибирский Федеральный Округ":         "Siberia",
		"Уральский Федеральный Округ":         "Ural",
		"Центральный Федеральный Округ":       "Central",
		"Южный Федеральный Округ":             "South",
	}

	query := `
		SELECT 
			"Federal_district",
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN "Quantity" END), 0) AS FOTON,
			COALESCE(SUM(CASE WHEN "Brand" = 'GAZ' THEN "Quantity" END), 0) AS GAZ,
			COALESCE(SUM(CASE WHEN "Brand" = 'ISUZU' THEN "Quantity" END), 0) AS ISUZU,
			COALESCE(SUM(CASE WHEN "Brand" = 'JAC' THEN "Quantity" END), 0) AS JAC,
			COALESCE(SUM(CASE WHEN "Brand" = 'KAMAZ' THEN "Quantity" END), 0) AS KAMAZ,
			COALESCE(SUM(CASE WHEN "Brand" NOT IN ('DONGFENG', 'FOTON', 'GAZ', 'ISUZU', 'JAC', 'KAMAZ') THEN "Quantity" END), 0) AS OTHER,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM ldt_3_5_12_truck_analytics_10_2024
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	dataByDistrict := orderedmap.New[string, []DistrictData]()
	customOrder := []string{
		"Summary",
		"Central",
		"North West",
		"Volga",
		"South",
		"North Caucasian",
		"Ural",
		"Siberia",
		"Far East",
	}

	for _, district := range customOrder {
		dataByDistrict.Set(district, []DistrictData{})
	}

	var summary DistrictData
	summary.RegionName = "Summary"
	initializeIntFields(&summary)

	for rows.Next() {
		var item DistrictData
		err := rows.Scan(
			&item.RegionName,
			&item.DONGFENG,
			&item.FOTON,
			&item.GAZ,
			&item.ISUZU,
			&item.JAC,
			&item.KAMAZ,
			&item.OTHER,
			&item.TotalMarket,
		)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}

		addToSummary(&summary, &item)

		translatedRegionName, exists := districtTranslations[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName
		}

		if existing, exists := dataByDistrict.Get(translatedRegionName); exists {
			item.RegionName = translatedRegionName
			dataByDistrict.Set(translatedRegionName, append(existing, item))
		}
	}

	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	if summaryData, exists := dataByDistrict.Get("Summary"); exists {
		dataByDistrict.Set("Summary", append(summaryData, summary))
	}

	ctx.JSON(http.StatusOK, Response{
		Data: dataByDistrict,
	})
}

func initializeIntFields(data *DistrictData) {
	data.DONGFENG = new(int)
	data.FOTON = new(int)
	data.GAZ = new(int)
	data.ISUZU = new(int)
	data.JAC = new(int)
	data.KAMAZ = new(int)
	data.OTHER = new(int)
}

func addToSummary(summary *DistrictData, item *DistrictData) {
	addIntFields(summary.DONGFENG, item.DONGFENG)
	addIntFields(summary.FOTON, item.FOTON)
	addIntFields(summary.GAZ, item.GAZ)
	addIntFields(summary.ISUZU, item.ISUZU)
	addIntFields(summary.JAC, item.JAC)
	addIntFields(summary.KAMAZ, item.KAMAZ)
	addIntFields(summary.OTHER, item.OTHER)
	summary.TotalMarket += item.TotalMarket
}

func addIntFields(a, b *int) {
	if a != nil && b != nil {
		*a += *b
	}
}
