const zone = window.location.pathname.slice("/zones/".length);

const red   = "rgb(178, 49, 49)";
const green = "rgb(19, 154, 19)";

const truckIconClassName = "fas fa-truck";
const humanIconClassName = "fas fa-street-view";
const inputIconClassName = "fa-solid fa-traffic-light";

let timeLeft    = 0;
let cameraIDs   = [];

window.onload = startPolling;

function disableButton(btn) {
    btn.disabled = true;
    btn.style.backgroundColor = red;
    btn.style.borderColor     = "rgb(94, 14, 14)";
    return btn
}

function enableButton(btn) {
    btn.disabled = false;
    btn.style.backgroundColor = green;
    btn.style.borderColor     = "rgb(10, 78, 10)";
    return btn
}

function changeColor(el, color) {
    el.style.color = color
}

async function get(url = "") {
    const response = await fetch(`${zone}/${url}`);
    return response.json();
}

async function startPolling() {
    const zone = newElement("p", {
        id: "zone",
        innerText: await get("/zone-name")
    });
    const countdownEl = updateCountdown(newElement("p", { id: "countdown" }));
    const statusBtnEl = disableButton(newElement("button", {
        id: "status-button",
        innerText: "Взвесить"
    }));
    const camerasDivEl = newElement("div", { id: "cams" });

    [zone, countdownEl, statusBtnEl, camerasDivEl].forEach(el => document.body.appendChild(el));
    statusBtnEl.onclick = () => console.log("click");

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
            enableButton(statusBtnEl);
            changeColor(countdownEl, green);
        } else {
            disableButton(statusBtnEl);
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
    return countdownEl
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
