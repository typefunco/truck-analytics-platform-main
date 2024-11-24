let form1 = document.querySelector(".filter__form-segment");
let form2 = document.querySelector(".filter__form-truck-type");
let form3 = document.querySelector(".filter__form-axle-type");
let form4 = document.querySelector(".filter__form-regions");
let form5 = document.querySelector(".filter__form-period");
let fotmBtn = document.querySelector(".filter-btn");
let btns1 = form1.querySelectorAll("input");
let btns2 = form2.querySelectorAll("input");
let btns3 = form3.querySelectorAll("input");
let btns4 = form4.querySelectorAll("input");
let btns5 = form5.querySelectorAll("span");
btns1.forEach((btn) => {
    btn.addEventListener("change", () => {
        form2.classList.add("activeForm");
    });
});
console.log(btns2);

btns2.forEach((btn) => {
    btn.addEventListener("change", () => {
        form3.classList.add("activeForm");
        let btns = document.querySelectorAll(".change-dis");
        btns.forEach((item) => {
            item.disabled = !item.disabled;
        });
    });
});
btns3.forEach((btn) => {
    btn.addEventListener("change", () => {
        form4.classList.add("activeForm");
        form5.classList.add("activeForm");
        fotmBtn.classList.add("activeForm");
    });
});

let spansPeriod = document.querySelectorAll(".period-span");

spansPeriod.forEach((span) => {
    span.addEventListener("click", () => {
        spansPeriod.forEach((item) => {
            // remove class
        });
        let currentSpan = span.dataset.period;
        for (let i = 0; i < currentSpan; i++) {
            // let temp = document.querySelector(
            //     '[data-period]="' + String(i) + '"'
            // );
            console.log('[data-period]="' + String(i) + '"');

            // document
            //     .querySelector('[data-period]="' + String(i) + '"')
            //     .classList.add("activeSpan");
        }
    });
});

const redirectBtn = document.querySelector(".filter-btn");
redirectBtn.addEventListener("click", (e) => {
    e.preventDefault();

    let truckType = document.querySelector('input[name="truck-type"]:checked')
        .dataset.value;
    let axleType = document.querySelector('input[name="axle-type"]:checked')
        .dataset.value;
    let regionsType = document.querySelector('input[name="regions"]:checked')
        .dataset.value;
    
    if (
        (truckType == "Tractors") &
        (axleType == "4x2") &
        (regionsType == "All regions")
    ) {
        window.location.href =
            "/analytics4x2Tractors/analytics4x2.html"; // Путь без /frontend
    } else if (
        (truckType == "Tractors") &
        (axleType == "4x2") &
        (regionsType == "Total Market")
    ) {
        window.location.href =
            "/analytics4x2TractorsTotalMarket/analytics4x2.html"; // Путь без /frontend
    } else if (
        (truckType == "Tractors") &
        (axleType == "6x4") &
        (regionsType == "Total Market")
    ) {
        window.location.href =
            "/analytics6x4TractorsTotalMarket/analytics6x4.html"; // Путь без /frontend
    } else if (
        (truckType == "Tractors") &
        (axleType == "6x4") &
        (regionsType == "All regions")
    ) {
        window.location.href =
            "/analytics6x4Tractors/analytics6x4.html"; // Путь без /frontend
    }

    if (
        (truckType == "Heavy chassis") &
        (axleType == "6x4") &
        (regionsType == "All regions")
    ) {
        window.location.href =
            "/analytics6x4Dumpers/analytics6x4.html"; // Путь без /frontend
    } else if (
        (truckType == "Heavy chassis") &
        (axleType == "6x4") &
        (regionsType == "Total Market")
    ) {
        window.location.href =
            "/analytics6x4DumpersTotalMarket/analytics6x4.html"; // Путь без /frontend
    } else if (
        (truckType == "Heavy chassis") &
        (axleType == "8x4") &
        (regionsType == "Total Market")
    ) {
        window.location.href =
            "/analytics8x4DumpersTotalMarket/analytics8x4.html"; // Путь без /frontend
    } else if (
        (truckType == "Heavy chassis") &
        (axleType == "8x4") &
        (regionsType == "All regions")
    ) {
        window.location.href =
            "/analytics8x4Dumpers/analytics8x4.html"; // Путь без /frontend
    }
});
