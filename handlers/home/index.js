const zone = window.location.pathname.slice("/zones/".length);

const countdownEl = document.getElementById("countdown");
const statusBtnEl = document.getElementById("status-button");
const camerasDivEl = document.getElementById("cams");

const red   = "rgb(178, 49, 49)";
const green = "rgb(19, 154, 19)";

const truckIconClassName = "fas fa-truck";
const humanIconClassName = "fas fa-street-view";
const inputIconClassName = "fa-solid fa-traffic-light";

let cameraIndex = 0;
let timeLeft    = 0;

window.onload       = startPolling;
statusBtnEl.onclick = () => console.log("click");

function disableStatusButton() {
    statusBtnEl.disabled = true;
    statusBtnEl.style.backgroundColor = red;
    statusBtnEl.style.borderColor     = "rgb(94, 14, 14)";
}

function enableStatusButton() {
    statusBtnEl.disabled = false;
    statusBtnEl.style.backgroundColor = green;
    statusBtnEl.style.borderColor     = "rgb(10, 78, 10)";
}

function changeColor(el, color) {
    el.style.color = color
}

async function get(url = "") {
    const response = await fetch(`${zone}/${url}`);
    return response.json();
}

function startPolling() {
    setTimeout(async function again() {
        const prevTimeLeft = timeLeft;
        timeLeft = await get("/countdown");
        if (timeLeft !== prevTimeLeft) updateCountdown(timeLeft);

        const cameraIDs = await get("/cameras-ids");
        for (const id of cameraIDs) {
            updateCameraStates(id, await get(`/cameras-info/${id}`));
        }

        if (!timeLeft) {
            enableStatusButton();
            changeColor(countdownEl, green);
        } else {
            disableStatusButton();
            changeColor(countdownEl, red);
        }

        setTimeout(again, 1000);
    });
}

function updateCountdown(timeLeft = 0) {
    const hours = formatNumber(Math.floor(timeLeft / 3600));
    const minutes = formatNumber(Math.floor(timeLeft / 60 - hours * 60));
    const seconds = formatNumber(timeLeft % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

formatNumber = (num) => num < 10 ? "0" + num: num;

function updateCameraStates(id, states) {
    const {car, human, inputs} = states;
    console.log(inputs)
    const camera = document.getElementById(id) ?? createCamera(id, inputs);
    for (const icon of camera.getElementsByTagName("i")) {
        if (icon.className.includes(truckIconClassName)) {
            setStatus(icon, car);
        } else if (icon.className.includes(humanIconClassName)) {
            setStatus(icon, human);
        } else if (icon.className.includes(inputIconClassName)) {
            setStatus(icon, Object.values(inputs).find(inp => inp.id === icon.id).state);
        }
    }
}

function createCamera(cameraID, inputs) {
    const camera = newElement("fieldset", { className: "cam", id: cameraID });
    const legend = newElement("legend", { innerText: `CAM-${++cameraIndex}` });
    const truckIcon = newElement("i", { className: truckIconClassName });
    const humanIcon = newElement("i", { className: humanIconClassName });
    const inputIcons = Object.values(inputs).map(inp => newElement("i", { className: inputIconClassName, id: inp.id }));

    [legend, truckIcon, humanIcon, ...inputIcons].forEach(el => camera.appendChild(el));
    camerasDivEl.appendChild(camera);
    return camera
}

function newElement(tagName, options = {}) {
    const el = document.createElement(tagName);
    for (const prop of Object.keys(options)) {
        el[prop] = options[prop];
    }

    return el
}

function setStatus(icon, status) {
    const isReady = icon.classList.contains("ready");
    if (isReady && !status || !isReady && status) {
        icon.classList.toggle("ready");
    }
}
