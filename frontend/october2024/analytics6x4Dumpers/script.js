document.addEventListener("DOMContentLoaded", () => {
    const urls = {
        2023: "http://localhost:8080/10m2023dumpers6x4",
        2024: "http://localhost:8080/10m2024dumpers6x4",
    };

    // Функция для загрузки данных
    async function fetchData(year) {
        try {
            const response = await fetch(urls[year]);
            if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);
            const jsonData = await response.json();

            if (year === 2023) {
                populateTable2023(jsonData); // Заполняем таблицу 2023 года
            } else if (year === 2024) {
                populateTable2024(jsonData); // Заполняем таблицу 2024 года
            }
        } catch (error) {
            console.error(`Error fetching data for ${year}:`, error);
            alert(
                `Failed to load data for ${year}. Check console for details.`
            );
        }
    }

    // Функция для создания графиков по регионам

    Chart.register(ChartDataLabels);

    function createRegionalChart(cell, regions, districtName) {
        // Заголовок округа
        const title = document.createElement("h3");
        title.textContent = `Federal region: ${districtName}`;
        title.style.color = "white";
        title.style.textAlign = "center";
        title.style.margin = "10px 0";
        cell.appendChild(title);

        // Контейнер для графиков
        const chartWrapper = document.createElement("div");
        chartWrapper.style.display = "flex"; // Размещение в ряд
        chartWrapper.style.justifyContent = "space-between";
        chartWrapper.style.marginBottom = "0px";
        cell.appendChild(chartWrapper);

        // Горизонтальный bar chart
        const barChartContainer = document.createElement("div");
        barChartContainer.style.flex = "1"; // Это гарантирует, что bar chart будет занимать доступное место
        chartWrapper.appendChild(barChartContainer);

        const margin = { top: 20, right: 150, bottom: 50, left: 150 };
        const width = 600 - margin.left - margin.right;
        const height = 400 - margin.top - margin.bottom;

        const svg = d3
            .select(barChartContainer)
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const brands = [
            "dongfeng",
            "faw",
            "howo",
            "jac",
            "sany",
            "sitrak",
            "shacman",
        ];
        const brandColors = {
            dongfeng: "#FF0000",
            faw: "#515a5a",
            shacman: "#FF6A00",
            howo: "#8B4513",
            jac: "#EA00FF",
            sany: "#4B0082",
            sitrak: "#00742C",
        };

        const filteredData = regions.filter(
            (region) => region.region_name !== districtName
        );

        const processedData = filteredData.map((d) => {
            const total = d.total || 1;
            return {
                region_name: d.region_name,
                total: d.total,
                ...brands.reduce((acc, brand) => {
                    acc[brand] = d[brand] || 0;
                    return acc;
                }, {}),
            };
        });

        const y = d3
            .scaleBand()
            .domain(processedData.map((d) => d.region_name))
            .range([0, height])
            .padding(0.1);

        const x = d3
            .scaleLinear()
            .domain([0, 1]) // Нормализованный масштаб
            .range([0, width]);

        // Добавляем название региона слева от каждого ряда
        svg.append("g")
            .selectAll("text.region-name")
            .data(processedData)
            .enter()
            .append("text")
            .attr("x", -margin.left / 2) // Сдвигаем текст влево от баров
            .attr("y", (d) => y(d.region_name) + y.bandwidth() / 2)
            .attr("dy", ".35em")
            .attr("fill", "white")
            .attr("text-anchor", "middle")
            .style("font-size", (d) =>
                d.region_name.length > 20 ? "8px" : "12px"
            ) // Уменьшаем шрифт для длинных названий
            .text((d) => d.region_name);

        svg.append("g")
            .selectAll("g")
            .data(processedData)
            .enter()
            .append("g")
            .attr("transform", (d) => `translate(0, ${y(d.region_name)})`)
            .selectAll("rect")
            .data((d) => {
                const total = d.total || 1;
                return brands.map((brand) => ({
                    brand,
                    value: d[brand] / total,
                    count: d[brand],
                    region: d.region_name,
                }));
            })
            .enter()
            .append("rect")
            .attr("x", (d, i, nodes) => {
                const previousValues = d3.sum(
                    brands.slice(0, brands.indexOf(d.brand)),
                    (brand) =>
                        d3.select(nodes[i].parentNode).data()[0][brand] /
                            d3.select(nodes[i].parentNode).data()[0].total || 0
                );
                return x(previousValues);
            })
            .attr("y", 0)
            .attr("width", (d) => x(d.value))
            .attr("height", y.bandwidth())
            .attr("fill", (d) => brandColors[d.brand])
            .append("title")
            .text((d) => `${d.region} - ${d.brand}: ${d.count} продаж`);

        svg.append("g")
            .selectAll("g")
            .data(processedData)
            .enter()
            .append("g")
            .attr("transform", (d) => `translate(0, ${y(d.region_name)})`)
            .selectAll("text")
            .data((d) => {
                const total = d.total || 1;
                return brands.map((brand) => ({
                    brand,
                    value: d[brand] / total,
                    count: d[brand],
                    region: d.region_name,
                }));
            })
            .enter()
            .append("text")
            .attr("x", (d, i, nodes) => {
                const previousValues = d3.sum(
                    brands.slice(0, brands.indexOf(d.brand)),
                    (brand) =>
                        d3.select(nodes[i].parentNode).data()[0][brand] /
                            d3.select(nodes[i].parentNode).data()[0].total || 0
                );
                return x(previousValues + d.value / 2);
            })
            .attr("y", y.bandwidth() / 2)
            .attr("dy", ".35em")
            .attr("fill", "white")
            .attr("text-anchor", "middle")
            .style("font-size", "12px")
            .text((d) => (d.count > 0 ? d.count : ""));

        // Добавление Pie Chart
        const totalMarketData = brands.map((brand) => {
            return {
                brand,
                value: d3.sum(filteredData, (d) => d[brand] || 0),
            };
        });

        const totalSales = d3.sum(totalMarketData, (d) => d.value);

        const pie = d3.pie().value((d) => d.value);
        const arc = d3.arc().innerRadius(0).outerRadius(150);

        const svgPie = d3
            .select(cell)
            .append("svg")
            .attr("width", 400)
            .attr("height", 400)
            .append("g")
            .attr("transform", `translate(200, 200)`);

        svgPie
            .selectAll("path")
            .data(pie(totalMarketData))
            .enter()
            .append("path")
            .attr("d", arc)
            .attr("fill", (d) => brandColors[d.data.brand])
            .attr("stroke", "white")
            .style("stroke-width", 1.5)
            .append("title") // Добавляем всплывающую подсказку
            .text((d) => {
                const percentage = (d.data.value / totalSales) * 100;
                return `${d.data.brand}: ${percentage.toFixed(1)}%`;
            });

        svgPie
            .selectAll("text")
            .data(pie(totalMarketData))
            .enter()
            .append("text")
            .attr("transform", (d) => `translate(${arc.centroid(d)})`)
            .attr("dy", ".35em")
            .attr("text-anchor", "middle")
            .attr("fill", "white")
            .style("font-size", "12px")
            .text((d) => {
                const percentage = (d.data.value / totalSales) * 100;
                return percentage >= 5 ? `${percentage.toFixed(1)}%` : "";
            });

        // Легенда
        const legendContainer = document.createElement("div");
        legendContainer.style.marginTop = "10px";
        legendContainer.style.color = "white";
        chartWrapper.appendChild(legendContainer);

        brands.forEach((brand) => {
            const legendItem = document.createElement("div");
            legendItem.style.display = "flex";
            legendItem.style.alignItems = "center";
            legendItem.style.marginBottom = "5px";

            const colorBox = document.createElement("div");
            colorBox.style.width = "20px";
            colorBox.style.height = "20px";
            colorBox.style.backgroundColor = brandColors[brand];
            colorBox.style.marginRight = "10px";

            const brandName = document.createElement("span");
            brandName.textContent =
                brand.charAt(0).toUpperCase() + brand.slice(1); // Преобразуем первую букву в верхний регистр

            legendItem.appendChild(colorBox);
            legendItem.appendChild(brandName);
            legendContainer.appendChild(legendItem);
        });
    }

    // Функция для заполнения таблицы за 2023 год
    function populateTable2023(data) {
        const tableBody = document.querySelector("#data-table-2023 tbody");
        tableBody.innerHTML = "";

        if (!data || !data.data) return;

        for (const district in data.data) {
            const regions = data.data[district];

            // Добавляем строку округа
            const districtRow = document.createElement("tr");
            districtRow.classList.add("district-row"); // Применяем класс для строки округа
            districtRow.innerHTML = `<td colspan="7"><strong>${district}</strong></td>`;
            tableBody.appendChild(districtRow);

            // Добавляем заголовок для регионов
            const brandRow = document.createElement("tr");
            brandRow.classList.add("brand-row"); // Применяем класс для строки заголовков
            brandRow.innerHTML = `
                <td><em>REGION</em></td>
                <td><em>FAW</em></td>
                <td><em>JAC</em></td>
                <td><em>SANY</em></td>
                <td><em>SITRAK</em></td>
                <td><em>TOTAL MARKET</em></td>
            `;
            tableBody.appendChild(brandRow);

            // Итерация по регионам
            regions.forEach((region) => {
                const row = document.createElement("tr");
                row.classList.add("data-row"); // Применяем класс для строки с данными
                row.innerHTML = `
                    <td>${region.region_name || "—"}</td>
                    <td>${region.faw || "—"}</td>
                    <td>${region.jac || "—"}</td>
                    <td>${region.sany || "—"}</td>
                    <td>${region.sitrak || "—"}</td>
                    <td>${region.total || "—"}</td>
                `;
                tableBody.appendChild(row);
            });

            // Добавляем строку для графика
            const chartRow = document.createElement("tr");
            const chartCell = document.createElement("td");
            chartCell.setAttribute("colspan", "7");
            chartCell.classList.add("pos"); // Применяем класс для стилизации ячейки графика
            chartRow.appendChild(chartCell);
            tableBody.appendChild(chartRow);

            // Создаем график для текущего округа
            createRegionalChart(chartCell, regions, district);
        }
    }

    // Функция для заполнения таблицы за 2024 год
    function populateTable2024(data) {
        const tableBody = document.querySelector("#data-table-2024 tbody");
        tableBody.innerHTML = "";

        if (!data || !data.data) return;

        for (const district in data.data) {
            const regions = data.data[district];

            // Добавляем строку округа
            const districtRow = document.createElement("tr");
            districtRow.classList.add("district-row"); // Применяем класс для оформления строки округа
            districtRow.innerHTML = `<td colspan="9"><strong>${district}</strong></td>`;
            tableBody.appendChild(districtRow);

            // Добавляем заголовок для регионов
            const brandRow = document.createElement("tr");
            brandRow.classList.add("brand-row"); // Применяем класс для строки заголовков
            brandRow.innerHTML = `
                <td><em>Region</em></td>
                <td><em>Faw</em></td>
                <td><em>Howo</em></td>
                <td><em>Jac</em></td>
                <td><em>Sany</em></td>
                <td><em>Sitrak</em></td>
                <td><em>Shacman</em></td>
                <td><em>Dongfeng</em></td>
                <td><em>Total Market</em></td>
            `;
            tableBody.appendChild(brandRow);

            // Итерация по регионам
            regions.forEach((region) => {
                const row = document.createElement("tr");
                row.classList.add("data-row"); // Применяем класс для строки с данными
                row.innerHTML = `
                    <td>${region.region_name || "—"}</td>
                    <td>${region.faw || "—"}</td>
                    <td>${region.howo || "—"}</td>
                    <td>${region.jac || "—"}</td>
                    <td>${region.sany || "—"}</td>
                    <td>${region.sitrak || "—"}</td>
                    <td>${region.shacman || "—"}</td>
                    <td>${region.dongfeng || "—"}</td>
                    <td>${region.total || "—"}</td>
                `;
                tableBody.appendChild(row);
            });

            // Добавляем строку для графика
            const chartRow = document.createElement("tr");
            const chartCell = document.createElement("td");
            chartCell.setAttribute("colspan", "9");
            chartCell.classList.add("pos"); // Применяем класс для стилизации ячейки графика
            chartRow.appendChild(chartCell);
            tableBody.appendChild(chartRow);

            // Создаем график для текущего округа
            createRegionalChart(chartCell, regions, district);
        }
    }

    // Загружаем данные для 2023 и 2024
    fetchData(2023);
    fetchData(2024);
});
