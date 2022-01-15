const zone = window.location.pathname.slice("/zones/".length);

const statusBtnEl = document.getElementById("status-button");

const red   = "rgb(178, 49, 49)";
const green = "rgb(19, 154, 19)";

const truckIconClassName = "fas fa-truck";
const humanIconClassName = "fas fa-street-view";
const inputIconClassName = "fa-solid fa-traffic-light";

let timeLeft    = 0;
let cameraIDs   = [];

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

async function startPolling() {
    const countdownEl = document.getElementById("countdown");
    const camerasDivEl = newElement("div", { id: "cams" });
    document.body.appendChild(camerasDivEl);

    document.getElementById("zone").innerText = await get("/zone-name");

    setTimeout(async function again() {
        const prevTimeLeft = timeLeft;
        timeLeft = await get("/countdown");
        if (timeLeft !== prevTimeLeft) updateCountdown(countdownEl, timeLeft);

        const prevCameraIDs = cameraIDs;
        cameraIDs = await get("/cameras-ids");

        const toRemove = prevCameraIDs.filter(prevID => !cameraIDs.find(currID => prevID === currID));
        toRemove.forEach(id => {
            const el = document.getElementById(id);
            el?.parentNode.removeChild(el);
        });

        for (const id of cameraIDs) {
            const camera = document.getElementById(id);
            const states = await get(`/cameras-info/${id}`)
            camera ? updateCameraStates(camera, states) : createCamera(camerasDivEl, id, states);
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

function updateCountdown(countdownEl, timeLeft = 0) {
    const hours = formatNumber(Math.floor(timeLeft / 3600));
    const minutes = formatNumber(Math.floor(timeLeft / 60 - hours * 60));
    const seconds = formatNumber(timeLeft % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

formatNumber = (num) => num < 10 ? "0" + num: num;

function createCamera(camerasDivEl, cameraID, states) {
    const {name, car, human, inputs} = states;
    const camera = newElement("fieldset", { className: "cam", id: cameraID });
    const legend = newElement("legend", { innerText: name });

    const truckIcon = newElement("i", { className: truckIconClassName });
    const humanIcon = newElement("i", { className: humanIconClassName });
    const inputIcons = Object.values(inputs).map(inp => newElement("i", {
        className: inputIconClassName,
        id: inp.id
    }));

    setStatus(truckIcon, car);
    setStatus(humanIcon, human);
    inputIcons.forEach(icon => setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state));

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

function updateCameraStates(camera, states) {
    const {car, human, inputs} = states;
    for (const icon of camera.getElementsByTagName("i")) {
        if (icon.className.includes(truckIconClassName)) {
            setStatus(icon, car);
        } else if (icon.className.includes(humanIconClassName)) {
            setStatus(icon, human);
        } else if (icon.className.includes(inputIconClassName)) {
            setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state);
        }
    }
}

function setStatus(icon, status) {
    const isReady = icon.classList.contains("ready");
    if (isReady && !status || !isReady && status) {
        icon.classList.toggle("ready");
    }
}
