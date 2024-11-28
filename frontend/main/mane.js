// Получаем элементы формы и сегменты
const form1 = document.querySelector(".filter__form-segment");
const form2 = document.querySelector(".filter__form-truck-type");
const form3 = document.querySelector(".filter__form-axle-type");
const form4 = document.querySelector(".filter__form-regions");
const form5 = document.querySelector(".filter__form-period");

// Находим все кнопки в сегменте
const segmentButtons = form1.querySelectorAll("input[name='segment']");

segmentButtons.forEach((btn) => {
    btn.addEventListener("change", (e) => {
        const selectedSegment = e.target.value; // Получаем выбранное значение сегмента
        console.log('Selected segment:', selectedSegment); // Логируем выбранное значение сегмента

        // Убираем все классы отображения
        form2.classList.add("hidden");
        form3.classList.add("hidden");
        form4.classList.add("hidden");
        form5.classList.add("hidden");

        if (selectedSegment === "HDT") {
            // Для HDT показываем все формы
            form2.classList.remove("hidden");
            form3.classList.remove("hidden");
            form4.classList.remove("hidden");
            form5.classList.remove("hidden");
        } else if (selectedSegment === "MDT" || selectedSegment === "LDT") {
            // Если выбран MDT или LDT, скрываем Truck Type и Axle Type
            form2.classList.add("hidden"); // Скрыть Truck Type
            form3.classList.add("hidden"); // Скрыть Axle Type

            // Убедиться, что Region и Period отображаются
            form4.classList.remove("hidden");
            form5.classList.remove("hidden");
        }
    });
});

// Логика кнопок Region и Period остается неизменной
const spansPeriod = document.querySelectorAll(".period-span");
spansPeriod.forEach((span) => {
    span.addEventListener("click", () => {
        spansPeriod.forEach((item) => {
            item.classList.remove("activeSpan"); // Убираем класс у всех
        });
        span.classList.add("activeSpan"); // Добавляем класс текущему
    });
});

const truckTypeInputs = document.querySelectorAll('input[name="truck-type"]');
const axleTypeInputs = document.querySelectorAll('input[name="axle-type"]');

truckTypeInputs.forEach((input) => {
    input.addEventListener("change", (e) => {
        const truckType = e.target.value;

        // Включаем/выключаем оси в зависимости от типа грузовика
        if (truckType === "Tractors") {
            // Если выбран "Tractors", то ось 8x4 заблокирована, а 4x2 и 6x4 доступны
            document.querySelector('input[value="8x4"]').disabled = true;  // Блокируем ось 8x4
            document.querySelector('input[value="4x2"]').disabled = false; // Разблокируем ось 4x2
            document.querySelector('input[value="6x4"]').disabled = false; // Разблокируем ось 6x4
        } else if (truckType === "Heavy chassis") {
            // Если выбран "Heavy chassis", то ось 4x2 заблокирована, а 6x4 и 8x4 доступны
            document.querySelector('input[value="8x4"]').disabled = false; // Разблокируем ось 8x4
            document.querySelector('input[value="4x2"]').disabled = true;  // Блокируем ось 4x2
            document.querySelector('input[value="6x4"]').disabled = false; // Разблокируем ось 6x4
        } else {
            // Если выбран другой тип (например, None), блокируем все оси
            document.querySelector('input[value="8x4"]').disabled = true;
            document.querySelector('input[value="4x2"]').disabled = true;
            document.querySelector('input[value="6x4"]').disabled = true;
        }
    });
});

// Логика перенаправления
const redirectBtn = document.querySelector(".filter-btn");
redirectBtn.addEventListener("click", (e) => {
    e.preventDefault();

    const truckTypeInput = document.querySelector('input[name="truck-type"]:checked');
    const axleTypeInput = document.querySelector('input[name="axle-type"]:checked');
    const regionsType = document.querySelector('input[name="regions"]:checked').dataset.value;

    const truckType = truckTypeInput ? truckTypeInput.dataset.value : null;
    const axleType = axleTypeInput ? axleTypeInput.dataset.value : null;

    // Логика перенаправления для MDT и LDT
    const selectedSegment = document.querySelector('input[name="segment"]:checked').value;

    if (
        selectedSegment === "MDT" &&
        regionsType === "Total Market"
    ) {
        window.location.href = `/mdtTotal/analytics4x2.html`;
    } else if (
        selectedSegment === "MDT" &&
        regionsType === "All regions"
    ) {
        window.location.href = `/mdt/analytics4x2.html`;
    } else if (
        selectedSegment === "LDT" &&
        regionsType === "Total Market")  {
        // Для LDT перенаправление
        window.location.href = `/ldtTotal/analytics4x2.html`;
    } else if (
        selectedSegment === "LDT" &&
        regionsType === "All regions") 
        {
        window.location.href = `/ldt/analytics4x2.html`;
    } else if (selectedSegment === "HDT") {
        // Логика для HDT (Heavy chassis)
        if (truckType === "Heavy chassis") {
            if (
                axleType === "6x4" &&
                regionsType === "All regions"
            ) {
                window.location.href = "/analytics6x4Dumpers/analytics6x4.html";
            } else if (
                axleType === "6x4" &&
                regionsType === "Total Market"
            ) {
                window.location.href = "/analytics6x4DumpersTotalMarket/analytics6x4.html";
            } else if (
                axleType === "8x4" &&
                regionsType === "All regions"
            ) {
                window.location.href = "/analytics8x4Dumpers/analytics8x4.html";
            } else if (
                axleType === "8x4" &&
                regionsType === "Total Market"
            ) {
                window.location.href = "/analytics8x4DumpersTotalMarket/analytics8x4.html";
            }
        } else if (truckType === "Tractors") {
            // Логика для тракторов
            if (
                axleType === "4x2" &&
                regionsType === "All regions"
            ) {
                window.location.href = "/analytics4x2Tractors/analytics4x2.html";
            } else if (
                axleType === "4x2" &&
                regionsType === "Total Market"
            ) {
                window.location.href = "/analytics4x2TractorsTotalMarket/analytics4x2.html";
            } else if (
                axleType === "6x4" &&
                regionsType === "Total Market"
            ) {
                window.location.href = "/analytics6x4TractorsTotalMarket/analytics6x4.html";
            } else if (
                axleType === "6x4" &&
                regionsType === "All regions"
            ) {
                window.location.href = "/analytics6x4Tractors/analytics6x4.html";
            }
        }
    }
});
