const countdownEl = document.getElementById("countdown");
const statusBtnEl = document.getElementById("status-button");
const camerasDivEl = document.getElementById("cams");

const truckIconClassName = "fas fa-truck";
const humanIconClassName = "fas fa-street-view";
const inputIconClassName = "fa-solid fa-traffic-light";

let cameraIndex = 0;
let timeLeft    = 0;

window.onload       = countdown;
statusBtnEl.onclick = async () => await recover();

function disableStatusButton() {
    statusBtnEl.disabled = true;
    statusBtnEl.style.backgroundColor = "rgb(178, 49, 49)";
    statusBtnEl.style.borderColor     = "rgb(94, 14, 14)";
}

function enableStatusButton() {
    statusBtnEl.disabled = false;
    statusBtnEl.style.backgroundColor = "rgb(19, 154, 19)";
    statusBtnEl.style.borderColor     = "rgb(10, 78, 10)";
}

async function get(url = "") {
    const response = await fetch(url);
    return response.json();
}

async function post(url = "", data = {}) {
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    return response.json();
}

async function recover() {
    try {
        await post("/reset-timer");
    } catch {
        console.log("timer have already reset");
    } finally {
        countdown();
        disableStatusButton();
    }
}

function countdown() {
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
            return
        }

        setTimeout(again, 1000);
    });
}

function updateCountdown(timeLeft = 0) {
    console.log("update countdown...")
    const hours = formatNumber(Math.floor(timeLeft / 3600));
    const minutes = formatNumber(Math.floor(timeLeft / 60 - hours * 60));
    const seconds = formatNumber(timeLeft % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

formatNumber = (num) => num < 10 ? "0" + num: num;

function updateCameraStates(id, states) {
    console.log("update camera states...")
    const {cars, humans, inputs} = states;
    const camera = document.getElementById(id) ?? createCamera(id);
    for (const icon of camera.getElementsByTagName("i")) {
        if (icon.className.includes(truckIconClassName)) {
            setStatus(icon, cars)
        } else if (icon.className.includes(humanIconClassName)) {
            setStatus(icon, humans)
        } else if (icon.className.includes(inputIconClassName)) {
            setStatus(icon, inputs)
        }
    }
}

function createCamera(cameraID) {
    const camera = newElement("fieldset", { className: "cam", id: cameraID });
    const legend = newElement("legend", { innerText: `CAM-${++cameraIndex}` });
    const truckIcon = newElement("i", { className: truckIconClassName });
    const humanIcon = newElement("i", { className: humanIconClassName })
    const inputIcon = newElement("i", { className: inputIconClassName });

    [legend, truckIcon, humanIcon, inputIcon].forEach(el => camera.appendChild(el));
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
