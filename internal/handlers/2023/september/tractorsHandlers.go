package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
)

func Home(ctx *gin.Context) {
	ctx.JSON(200, "Cool")
}

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
	"Центральный Федеральный Округ":       "Central Federal District",
	"Северо-Западный Федеральный Округ":   "North West Federal District",
	"Южный Федеральный Округ":             "South Federal District",
	"Северо-Кавказский Федеральный Округ": "North Caucasian Federal District",
	"Приволжский Федеральный Округ":       "Volga Federal District",
	"Уральский Федеральный Округ":         "Ural Federal District",
	"Сибирский Федеральный Округ":         "Siberia Federal District",
	"Дальневосточный Федеральный Округ":   "Far East Federal District",
}

// Мапа с переводами федеральных округов на английский
var districtTranslations = map[string]string{
	"Центральный Федеральный Округ":       "Central Federal District",
	"Северо-Западный Федеральный Округ":   "North West Federal District",
	"Южный Федеральный Округ":             "South Federal District",
	"Северо-Кавказский Федеральный Округ": "North Caucasian Federal District",
	"Приволжский Федеральный Округ":       "Volga Federal District",
	"Уральский Федеральный Округ":         "Ural Federal District",
	"Сибирский Федеральный Округ":         "Siberia Federal District",
	"Дальневосточный Федеральный Округ":   "Far East Federal District",
}

func NineMonth2023Tractors4x2(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName string `json:"region_name"`
		DONGFENG   *int   `json:"dongfeng"`
		FAW        *int   `json:"faw"`
		FOTON      *int   `json:"foton"`
		JAC        *int   `json:"jac"`
		SHACMAN    *int   `json:"shacman"`
		SITRAK     *int   `json:"sitrak"`
		TOTAL      int    `json:"total"`
	}

	type TruckAnalyticsResponse struct {
		Data  map[string][]TruckAnalytics `json:"data"`
		Error string                      `json:"error,omitempty"`
	}

	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2023_01_12."Federal_district",
				truck_analytics_2023_01_12."Region",
				truck_analytics_2023_01_12."Brand",
				SUM(truck_analytics_2023_01_12."Quantity") as total_sales
			FROM truck_analytics_2023_01_12
			WHERE 
				truck_analytics_2023_01_12."Wheel_formula" = '4x2'
				AND truck_analytics_2023_01_12."Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'JAC', 'SHACMAN', 'SITRAK')
				AND truck_analytics_2023_01_12."Month_of_registration" <= 9
				AND truck_analytics_2023_01_12."Body_type" = 'Седельный тягач'
				AND truck_analytics_2023_01_12."Exact_mass" = 18000
			GROUP BY 
				truck_analytics_2023_01_12."Federal_district", 
				truck_analytics_2023_01_12."Region", 
				truck_analytics_2023_01_12."Brand"
		)
		SELECT 
			"Federal_district",
			COALESCE("Region", "Federal_district") as Region_name,
			MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) as DONGFENG,
			MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
			MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) as FOTON,
			MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) as JAC,
			MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END), 0) as TOTAL
		FROM base_data
		GROUP BY 
			"Federal_district",
			"Region"
		ORDER BY 
			"Federal_district",
			CASE 
				WHEN "Region" = "Federal_district" THEN 1 
				ELSE 0 
			END,
			"Region"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Failed to execute query: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	dataByDistrict := make(map[string][]TruckAnalytics)
	summaryByDistrict := make(map[string]TruckAnalytics)

	// Читаем строки результатов
	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.DONGFENG,
			&ta.FAW,
			&ta.FOTON,
			&ta.JAC,
			&ta.SHACMAN,
			&ta.SITRAK,
			&ta.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
				Error: "Failed to scan row: " + err.Error(),
			})
			return
		}

		// Перевод названий
		if engName, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = engName
		}
		if engDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = engDistrict
		}

		// Обновляем суммарные данные по округу
		summary := summaryByDistrict[federalDistrict]

		if ta.DONGFENG != nil {
			if summary.DONGFENG == nil {
				summary.DONGFENG = new(int)
			}
			*summary.DONGFENG += *ta.DONGFENG
		}
		if ta.FAW != nil {
			if summary.FAW == nil {
				summary.FAW = new(int)
			}
			*summary.FAW += *ta.FAW
		}
		if ta.FOTON != nil {
			if summary.FOTON == nil {
				summary.FOTON = new(int)
			}
			*summary.FOTON += *ta.FOTON
		}
		if ta.JAC != nil {
			if summary.JAC == nil {
				summary.JAC = new(int)
			}
			*summary.JAC += *ta.JAC
		}
		if ta.SHACMAN != nil {
			if summary.SHACMAN == nil {
				summary.SHACMAN = new(int)
			}
			*summary.SHACMAN += *ta.SHACMAN
		}
		if ta.SITRAK != nil {
			if summary.SITRAK == nil {
				summary.SITRAK = new(int)
			}
			*summary.SITRAK += *ta.SITRAK
		}

		summary.TOTAL += ta.TOTAL
		summary.RegionName = federalDistrict
		summaryByDistrict[federalDistrict] = summary

		// Добавляем данные региона в соответствующий федеральный округ
		dataByDistrict[federalDistrict] = append(dataByDistrict[federalDistrict], ta)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	// Добавляем суммарные данные по федеральному округу
	for district, summary := range summaryByDistrict {
		dataByDistrict[district] = append(dataByDistrict[district], summary)
	}

	// Отправляем ответ
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

func Tractors4x2WithTotalMarket2023(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		RegionName string `json:"region_name"`
		Dongfeng   *int   `json:"dongfeng"`
		Faw        *int   `json:"faw"`
		Foton      *int   `json:"foton"`
		Jac        *int   `json:"jac"`
		Shacman    *int   `json:"shacman"`
		Sitrak     *int   `json:"sitrak"`
		Total      int    `json:"total"`
	}

	// Структура для итогов
	type Summary struct {
		Dongfeng int `json:"dongfeng"`
		Faw      int `json:"faw"`
		Foton    int `json:"foton"`
		Jac      int `json:"jac"`
		Shacman  int `json:"shacman"`
		Sitrak   int `json:"sitrak"`
		Total    int `json:"total"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far Eastern Federal District",
		"Приволжский Федеральный Округ":       "Volga Federal District",
		"Северо-Западный Федеральный Округ":   "North West Federal District",
		"Северо-Кавказский Федеральный Округ": "North Caucasian Federal District",
		"Сибирский Федеральный Округ":         "Siberian Federal District",
		"Уральский Федеральный Округ":         "Ural Federal District",
		"Центральный Федеральный Округ":       "Central Federal District",
		"Южный Федеральный Округ":             "Southern Federal District",
	}

	// SQL-запрос для получения данных по округам
	query := `
		SELECT 
			"Federal_district",
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM(CASE WHEN "Brand" = 'FAW' THEN "Quantity" END), 0) AS FAW,
			COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN "Quantity" END), 0) AS FOTON,
			COALESCE(SUM(CASE WHEN "Brand" = 'JAC' THEN "Quantity" END), 0) AS JAC,
			COALESCE(SUM(CASE WHEN "Brand" = 'SHACMAN' THEN "Quantity" END), 0) AS SHACMAN,
			COALESCE(SUM(CASE WHEN "Brand" = 'SITRAK' THEN "Quantity" END), 0) AS SITRAK,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2023_01_12
		WHERE 
			"Wheel_formula" = '4x2'
			AND "Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'JAC', 'SHACMAN', 'SITRAK')
			AND "Month_of_registration" <= 9
			AND "Body_type" = 'Седельный тягач'
			AND "Exact_mass" = 18000
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	// Подключение к базе данных
	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	// Запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	var data []DistrictData
	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Dongfeng, &item.Faw, &item.Foton, &item.Jac, &item.Shacman, &item.Sitrak, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		data = append(data, item)
	}
	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Итоговые данные
	var summary Summary
	for _, item := range data {
		if item.Dongfeng != nil {
			summary.Dongfeng += *item.Dongfeng
		}
		if item.Faw != nil {
			summary.Faw += *item.Faw
		}
		if item.Foton != nil {
			summary.Foton += *item.Foton
		}
		if item.Jac != nil {
			summary.Jac += *item.Jac
		}
		if item.Shacman != nil {
			summary.Shacman += *item.Shacman
		}
		if item.Sitrak != nil {
			summary.Sitrak += *item.Sitrak
		}
		summary.Total += item.Total
	}

	// Подготовка JSON-ответа
	response := make(map[string][]DistrictData)

	// Добавляем Summary в начало
	response["Summary"] = []DistrictData{
		{
			RegionName: "Summary",
			Dongfeng:   &summary.Dongfeng,
			Faw:        &summary.Faw,
			Foton:      &summary.Foton,
			Jac:        &summary.Jac,
			Shacman:    &summary.Shacman,
			Sitrak:     &summary.Sitrak,
			Total:      summary.Total,
		},
	}

	// Добавляем данные по округам после Summary
	for _, item := range data {
		// Перевод имени округа
		translatedRegionName, exists := regionTranslation[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName // если нет перевода, оставляем как есть
		}

		// Вставляем данные для каждого региона
		if _, exists := response[translatedRegionName]; !exists {
			response[translatedRegionName] = []DistrictData{}
		}
		item.RegionName = translatedRegionName
		response[translatedRegionName] = append(response[translatedRegionName], item)
	}

	// Ответ
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func NineMonth2023Tractors6x4(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName string `json:"region_name"`
		DONGFENG   *int   `json:"dongfeng"`
		FAW        *int   `json:"faw"`
		FOTON      *int   `json:"foton"`
		HOWO       *int   `json:"howo"`
		SHACMAN    *int   `json:"shacman"`
		SITRAK     *int   `json:"sitrak"`
		TOTAL      int    `json:"total"`
	}

	type TruckAnalyticsResponse struct {
		Data  map[string][]TruckAnalytics `json:"data"`
		Error string                      `json:"error,omitempty"`
	}

	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2023_01_12."Federal_district",
				truck_analytics_2023_01_12."Region",
				truck_analytics_2023_01_12."Brand",
				SUM(truck_analytics_2023_01_12."Quantity") as total_sales
			FROM truck_analytics_2023_01_12
			WHERE 
				truck_analytics_2023_01_12."Wheel_formula" = '6x4'
				AND truck_analytics_2023_01_12."Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'HOWO', 'SHACMAN', 'SITRAK')
				AND truck_analytics_2023_01_12."Month_of_registration" <= 9
				AND truck_analytics_2023_01_12."Body_type" = 'Седельный тягач'
				AND truck_analytics_2023_01_12."Exact_mass" = 25000
			GROUP BY 
				truck_analytics_2023_01_12."Federal_district", 
				truck_analytics_2023_01_12."Region", 
				truck_analytics_2023_01_12."Brand"
		),
		federal_totals AS (
			SELECT 
				"Federal_district",
				"Federal_district" as "Region",
				"Brand",
				SUM(total_sales) as total_sales
			FROM base_data
			GROUP BY "Federal_district", "Brand"
		),
		combined_data AS (
			SELECT * FROM base_data
			UNION ALL
			SELECT * FROM federal_totals
		)
		SELECT 
			"Federal_district",
			COALESCE("Region", "Federal_district") as Region_name,
			MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) as DONGFENG,
			MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
			MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) as FOTON,
			MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END) as HOWO,
			MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END), 0) as TOTAL
		FROM combined_data
		GROUP BY 
			"Federal_district",
			"Region"
		ORDER BY 
			"Federal_district",
			CASE 
				WHEN "Region" = "Federal_district" THEN 1 
				ELSE 0 
			END,
			"Region"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Failed to execute query: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	dataByDistrict := make(map[string][]TruckAnalytics)

	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.DONGFENG,
			&ta.FAW,
			&ta.FOTON,
			&ta.HOWO,
			&ta.SHACMAN,
			&ta.SITRAK,
			&ta.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
				Error: "Failed to scan row: " + err.Error(),
			})
			return
		}

		// Переводим названия на английский
		if engName, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = engName
		}
		if engDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = engDistrict
		}

		// Добавляем данные региона
		dataByDistrict[federalDistrict] = append(dataByDistrict[federalDistrict], ta)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	// Добавляем суммарные данные по федеральным округам
	for district, regions := range dataByDistrict {
		totalDongfeng, totalFaw, totalFoton, totalHowo := 0, 0, 0, 0
		totalShacman, totalSitrak, total := 0, 0, 0

		for _, region := range regions {
			if region.DONGFENG != nil {
				totalDongfeng += *region.DONGFENG
			}
			if region.FAW != nil {
				totalFaw += *region.FAW
			}
			if region.FOTON != nil {
				totalFoton += *region.FOTON
			}
			if region.HOWO != nil {
				totalHowo += *region.HOWO
			}
			if region.SHACMAN != nil {
				totalShacman += *region.SHACMAN
			}
			if region.SITRAK != nil {
				totalSitrak += *region.SITRAK
			}
			total += region.TOTAL
		}

		// Проверяем, если федеральный округ уже есть в данных, не добавляем его снова
		if _, exists := dataByDistrict[district]; !exists {
			dataByDistrict[district] = append(dataByDistrict[district], TruckAnalytics{
				RegionName: district,
				DONGFENG:   &totalDongfeng,
				FAW:        &totalFaw,
				FOTON:      &totalFoton,
				HOWO:       &totalHowo,
				SHACMAN:    &totalShacman,
				SITRAK:     &totalSitrak,
				TOTAL:      total,
			})
		}
	}

	// Отправляем ответ
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

func Tractors6x4WithTotalMarket(ctx *gin.Context) {
	type BrandAnalytics struct {
		DistrictName string `json:"district_name"`
		DONGFENG     *int   `json:"dongfeng"`
		FAW          *int   `json:"faw"`
		FOTON        *int   `json:"foton"`
		HOWO         *int   `json:"howo"`
		SHACMAN      *int   `json:"shacman"`
		SITRAK       *int   `json:"sitrak"`
		TOTAL        int    `json:"total"`
	}

	type FederalDistrictAnalyticsResponse struct {
		Data  map[string][]BrandAnalytics `json:"data"`
		Total int                         `json:"total_market"`
		Error string                      `json:"error,omitempty"`
	}

	query := `
		WITH base_data AS (
			SELECT 
				truck_analytics_2023_01_12."Federal_district",
				truck_analytics_2023_01_12."Brand",
				SUM(truck_analytics_2023_01_12."Quantity") as total_sales
			FROM truck_analytics_2023_01_12
			WHERE 
				truck_analytics_2023_01_12."Wheel_formula" = '6x4'
				AND truck_analytics_2023_01_12."Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'HOWO', 'SHACMAN', 'SITRAK')
				AND truck_analytics_2023_01_12."Month_of_registration" <= 9
				AND truck_analytics_2023_01_12."Body_type" = 'Седельный тягач'
				AND truck_analytics_2023_01_12."Exact_mass" = 25000
			GROUP BY 
				truck_analytics_2023_01_12."Federal_district", 
				truck_analytics_2023_01_12."Brand"
		)
		SELECT 
			"Federal_district",
			MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) as DONGFENG,
			MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END) as FAW,
			MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) as FOTON,
			MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END) as HOWO,
			MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END) as SHACMAN,
			MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END) as SITRAK,
			COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FAW' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SHACMAN' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'SITRAK' THEN total_sales END), 0) as TOTAL
		FROM base_data
		GROUP BY 
			"Federal_district"
		ORDER BY 
			"Federal_district"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Warn("Can't connect to database")
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, FederalDistrictAnalyticsResponse{
			Error: "Failed to execute query: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	dataByDistrict := make(map[string][]BrandAnalytics)
	totalMarket := 0

	// Читаем строки результатов
	for rows.Next() {
		var ba BrandAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ba.DONGFENG,
			&ba.FAW,
			&ba.FOTON,
			&ba.HOWO,
			&ba.SHACMAN,
			&ba.SITRAK,
			&ba.TOTAL,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, FederalDistrictAnalyticsResponse{
				Error: "Failed to scan row: " + err.Error(),
			})
			return
		}

		// Перевод названия округа
		if engDistrict, ok := districtTranslations[federalDistrict]; ok {
			ba.DistrictName = engDistrict
		} else {
			ba.DistrictName = federalDistrict
		}

		// Суммируем общий рынок
		totalMarket += ba.TOTAL

		// Группируем данные по федеральным округам
		dataByDistrict[federalDistrict] = append(dataByDistrict[federalDistrict], ba)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, FederalDistrictAnalyticsResponse{
			Error: "Error iterating over rows: " + err.Error(),
		})
		return
	}

	// Отправляем ответ
	ctx.JSON(http.StatusOK, FederalDistrictAnalyticsResponse{
		Data:  dataByDistrict,
		Total: totalMarket,
	})
}

func Tractors6x4WithTotalMarket2023(ctx *gin.Context) {
	// Структура для хранения данных округа
	type DistrictData struct {
		RegionName string `json:"region_name"`
		Dongfeng   *int   `json:"dongfeng"`
		Faw        *int   `json:"faw"`
		Foton      *int   `json:"foton"`
		Howo       *int   `json:"howo"`
		Shacman    *int   `json:"shacman"`
		Sitrak     *int   `json:"sitrak"`
		Total      int    `json:"total"`
	}

	// Структура для итогов
	type Summary struct {
		Dongfeng int `json:"dongfeng"`
		Faw      int `json:"faw"`
		Foton    int `json:"foton"`
		Howo     int `json:"howo"`
		Shacman  int `json:"shacman"`
		Sitrak   int `json:"sitrak"`
		Total    int `json:"total"`
	}

	// Мапа для перевода русских названий федеральных округов на английский
	var regionTranslation = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far Eastern Federal District",
		"Приволжский Федеральный Округ":       "Volga Federal District",
		"Северо-Западный Федеральный Округ":   "North West Federal District",
		"Северо-Кавказский Федеральный Округ": "North Caucasian Federal District",
		"Сибирский Федеральный Округ":         "Siberian Federal District",
		"Уральский Федеральный Округ":         "Ural Federal District",
		"Центральный Федеральный Округ":       "Central Federal District",
		"Южный Федеральный Округ":             "Southern Federal District",
	}

	// SQL-запрос для получения данных по округам
	query := `
		SELECT 
			"Federal_district",
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM(CASE WHEN "Brand" = 'FAW' THEN "Quantity" END), 0) AS FAW,
			COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN "Quantity" END), 0) AS FOTON,
			COALESCE(SUM(CASE WHEN "Brand" = 'HOWO' THEN "Quantity" END), 0) AS HOWO,
			COALESCE(SUM(CASE WHEN "Brand" = 'SHACMAN' THEN "Quantity" END), 0) AS SHACMAN,
			COALESCE(SUM(CASE WHEN "Brand" = 'SITRAK' THEN "Quantity" END), 0) AS SITRAK,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM truck_analytics_2023_01_12
		WHERE 
			"Wheel_formula" = '6x4'
			AND "Brand" IN ('DONGFENG', 'FAW', 'FOTON', 'HOWO', 'SHACMAN', 'SITRAK')
			AND "Month_of_registration" <= 9
			AND "Body_type" = 'Седельный тягач'
			AND "Exact_mass" = 25000
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	// Подключение к базе данных
	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	// Запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	var data []DistrictData
	for rows.Next() {
		var item DistrictData
		err := rows.Scan(&item.RegionName, &item.Dongfeng, &item.Faw, &item.Foton, &item.Howo, &item.Shacman, &item.Sitrak, &item.Total)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		data = append(data, item)
	}
	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	// Итоговые данные
	var summary Summary
	for _, item := range data {
		if item.Dongfeng != nil {
			summary.Dongfeng += *item.Dongfeng
		}
		if item.Faw != nil {
			summary.Faw += *item.Faw
		}
		if item.Foton != nil {
			summary.Foton += *item.Foton
		}
		if item.Howo != nil {
			summary.Howo += *item.Howo
		}
		if item.Shacman != nil {
			summary.Shacman += *item.Shacman
		}
		if item.Sitrak != nil {
			summary.Sitrak += *item.Sitrak
		}
		summary.Total += item.Total
	}

	// Подготовка JSON-ответа
	response := make(map[string][]DistrictData)

	// Добавляем Summary в начало
	response["Summary"] = []DistrictData{
		{
			RegionName: "Summary",
			Dongfeng:   &summary.Dongfeng,
			Faw:        &summary.Faw,
			Foton:      &summary.Foton,
			Howo:       &summary.Howo,
			Shacman:    &summary.Shacman,
			Sitrak:     &summary.Sitrak,
			Total:      summary.Total,
		},
	}

	// Добавляем данные по округам после Summary
	for _, item := range data {
		// Перевод имени округа
		translatedRegionName, exists := regionTranslation[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName // если нет перевода, оставляем как есть
		}

		// Вставляем данные для каждого региона
		if _, exists := response[translatedRegionName]; !exists {
			response[translatedRegionName] = []DistrictData{}
		}
		item.RegionName = translatedRegionName
		response[translatedRegionName] = append(response[translatedRegionName], item)
	}

	// Ответ
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
