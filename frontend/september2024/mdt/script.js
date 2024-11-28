document.addEventListener("DOMContentLoaded", () => {
    const urls = {
        2023: "http://localhost:8080/9m2023tractors4x2",
        2024: "http://localhost:8080/9m2024tractors4x2",
    };

    async function fetchData(year) {
        try {
            const response = await fetch(urls[year]);
            if (!response.ok) throw new Error(`HTTP Error: ${response.status}`);
            const jsonData = await response.json();
            populateTable(jsonData, year);
        } catch (error) {
            console.error(`Error fetching data for ${year}:`, error);
            alert(
                `Failed to load data for ${year}. Check console for details.`
            );
        }
    }


    // Функция для создания горизонтального графика по регионам
    // Подключаем DataLabels плагин для текста на барах
    Chart.register(ChartDataLabels);

    // Функция для создания графиков
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
        chartWrapper.style.display = "flex";  // Размещение в ряд
        chartWrapper.style.justifyContent = "space-between";
        chartWrapper.style.marginBottom = "0px";
        cell.appendChild(chartWrapper);
    
        // Горизонтальный bar chart
        const barChartContainer = document.createElement("div");
        barChartContainer.style.flex = "1"; // Это гарантирует, что bar chart будет занимать доступное место
        chartWrapper.appendChild(barChartContainer);
    
        const margin = { top: 20, right: 150, bottom: 50, left: 150 };
        const width = 800 - margin.left - margin.right;
        const height = 400 - margin.top - margin.bottom;
    
        const svg = d3.select(barChartContainer)
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);
    
        // Список брендов и цвета
        const brands = ["dongfeng", "faw", "foton", "jac", "shacman", "sitrak"];
        const brandColors = {
            dongfeng: "#FF0000",
            faw: "#999999",
            foton: "#002BFF",
            jac: "#EA00FF",
            shacman: "#FF6A00",
            sitrak: "#00742C"
        };
    
        // Фильтрация данных
        const filteredData = regions.filter((region) => region.region_name !== districtName);
    
        // Обработка данных
        const processedData = filteredData.map((d) => {
            const total = d.total || 1;
            return {
                region_name: d.region_name,
                total: d.total,
                ...brands.reduce((acc, brand) => {
                    acc[brand] = d[brand] || 0;
                    return acc;
                }, {})
            };
        });
    
        // Масштабы
        const y = d3.scaleBand()
            .domain(processedData.map((d) => d.region_name))
            .range([0, height])
            .padding(0.1);
    
        const x = d3.scaleLinear()
            .domain([0, 1]) // Нормализованный масштаб
            .range([0, width]);
    
        // Отрисовка stacked bar chart
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
                    region: d.region_name
                }));
            })
            .enter()
            .append("rect")
            .attr("x", (d, i, nodes) => {
                const previousValues = d3.sum(
                    brands.slice(0, brands.indexOf(d.brand)),
                    (brand) => d3.select(nodes[i].parentNode).data()[0][brand] / d3.select(nodes[i].parentNode).data()[0].total || 0
                );
                return x(previousValues);
            })
            .attr("y", 0)
            .attr("width", (d) => x(d.value))
            .attr("height", y.bandwidth())
            .attr("fill", (d) => brandColors[d.brand])
            .append("title")
            .text((d) => `${d.region} - ${d.brand}: ${d.count} продаж`);
    
        // Добавление текста с названиями регионов для bar chart
        svg.append("g")
            .selectAll("g")
            .data(processedData)
            .enter()
            .append("g")
            .attr("transform", (d) => `translate(0, ${y(d.region_name)})`)
            .append("text")
            .attr("x", -5)  // Сдвигаем текст слева от бара
            .attr("y", y.bandwidth() / 2)
            .attr("dy", ".35em")
            .attr("fill", "white")
            .style("text-anchor", "end")
            .style("font-size", (d) => d.region_name.length > 20 ? "8px" : "12px")  // Уменьшаем шрифт для длинных названий
            .text((d) => d.region_name);
    
        // Добавление текста с продажами, только если значение > 5%
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
                    region: d.region_name
                }));
            })
            .enter()
            .append("text")
            .attr("x", (d, i, nodes) => {
                const previousValues = d3.sum(
                    brands.slice(0, brands.indexOf(d.brand)),
                    (brand) => d3.select(nodes[i].parentNode).data()[0][brand] / d3.select(nodes[i].parentNode).data()[0].total || 0
                );
                return x(previousValues + d.value / 2);
            })
            .attr("y", y.bandwidth() / 2)
            .attr("dy", ".35em")
            .attr("fill", "white")
            .attr("text-anchor", "middle")
            .style("font-size", "12px")
            .text((d) => {
                const total = d.total || 1;
                const percentage = (d.count / total) * 100;
                return percentage >= 5 ? d.count : ""; // Показываем только если значение > 5%
            });
    
        // Добавление легенды для bar chart
        const legend = svg.append("g")
            .attr("transform", `translate(${width + 20}, 0)`);
    
        brands.forEach((brand, i) => {
            legend.append("rect")
                .attr("x", 0)
                .attr("y", i * 20)
                .attr("width", 18)
                .attr("height", 18)
                .style("fill", brandColors[brand]);
    
            legend.append("text")
                .attr("x", 24)
                .attr("y", i * 20 + 9)
                .attr("dy", ".35em")
                .style("fill", "white")
                .text(brand);
        });
    
        // Отрисовка Pie Chart
        const totalMarketData = brands.map((brand) => {
            return {
                brand,
                value: d3.sum(filteredData, (d) => d[brand] || 0)
            };
        });
    
        const totalSales = d3.sum(totalMarketData, (d) => d.value);
    
        const pie = d3.pie().value((d) => d.value);
        const arc = d3.arc().innerRadius(0).outerRadius(150);
    
        const svgPie = d3.select(cell).append("svg")
            .attr("width", 400)
            .attr("height", 400)
            .append("g")
            .attr("transform", `translate(200, 200)`);
    
        // Отрисовка сегментов на Pie Chart
        svgPie.selectAll("path")
            .data(pie(totalMarketData))
            .enter()
            .append("path")
            .attr("d", arc)
            .attr("fill", (d) => brandColors[d.data.brand])
            .attr("stroke", "white")
            .style("stroke-width", 1.5);
    
        // Отображение процентов на Pie Chart
        svgPie.selectAll("text")
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
    
    }
    
    

    // Обновляем populateTable для графиков по регионам
    function populateTable(data, year) {
        const tableBody = document.querySelector(`#data-table-${year} tbody`);
        tableBody.innerHTML = "";
    
        if (!data || !data.data) return;
    
        const districts = Object.entries(data.data).slice(1);
    
        for (const [district, regions] of districts) {
            const districtRow = document.createElement("tr");
            districtRow.classList.add("district-row");
            districtRow.innerHTML = `<td colspan="8"><strong>${district}</strong></td>`;
            tableBody.appendChild(districtRow);
    
            const brandRow = document.createElement("tr");
            brandRow.classList.add("brand-row");
            brandRow.innerHTML = `
                <td><em>Region</em></td>
                <td><em>Foton</em></td>
                <td><em>Faw</em></td>
                <td><em>Dongfeng</em></td>
                <td><em>Jac</em></td>
                <td><em>Shacman</em></td>
                <td><em>Sitrak</em></td>
                <td><em>Total Market</em></td>
            `;
            tableBody.appendChild(brandRow);
    
            regions.forEach((region) => {
                const row = document.createElement("tr");
                row.classList.add("data-row");
                row.innerHTML = `
                    <td>${region.region_name || "—"}</td>
                    <td>${region.foton ?? "—"}</td>
                    <td>${region.faw ?? "—"}</td>
                    <td>${region.dongfeng || "—"}</td>
                    <td>${region.jac ?? "—"}</td>
                    <td>${region.shacman ?? "—"}</td>
                    <td>${region.sitrak ?? "—"}</td>
                    <td>${region.total ?? "—"}</td>
                `;
                tableBody.appendChild(row);
            });
    
            const chartRow = document.createElement("tr");
            const chartCell = document.createElement("td");
            chartCell.setAttribute("colspan", "8");
            chartCell.classList.add("pos");
            chartRow.appendChild(chartCell);
            tableBody.appendChild(chartRow);
    
            createRegionalChart(chartCell, regions, district);
        }
        console.log('Styles applied to district and brand rows'); // добавьте эту строку для отладки
    }
    
    

    // Вызываем fetchData для каждого года
    fetchData("2023");
    fetchData("2024");
});