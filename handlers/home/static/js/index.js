const zone = window.location.pathname.slice("/zones/".length);

const pollingFrequency = 1000; // in ms

const truckIconClassName = "icon truck";
const humanIconClassName = "icon human";
const inputIconClassName = "icon input";

const linkSlashClassName = "icon link-slash"
const heartCrackClassName = "icon heart-crack";

const circleCheckClassName = "icon circle-check"
const circleXmarkClassName = "icon circle-xmark";

async function get(url = "", format = "json") {
    const resp = await fetch(`${zone}/${url}`);
    if (!resp.ok) throw new FetchError(await resp.text());

    return format.toLowerCase() === "text" ? resp.text() : resp.json();
}

class FetchError extends Error {
    constructor(text) {
        super(`error: ${text}`);
        this.name = this.constructor.name;
    }
}

function newElement(tagName, options = {}) {
    const el = document.createElement(tagName);
    for (const prop of Object.keys(options)) {
        el[prop] = options[prop];
    }

    return el;
}

class App {
    constructor() {
        this.cameras = {};
        this.isHealthy = false;
        this.bodyEl = document.getElementById("body-container");
        this.spinnerEl = document.getElementById("spinner");
    }

    render(name = "") {
        this.bodyEl.innerHTML = `
            <p id="zone">${name}</p>
            <p id="countdown">00:00:00</p>
            <button id="status-button" disabled>Зважити</button>
            <fieldset id="status-button-error" style="display: none">
                <p class="error"></p>
            </fieldset>
            <div id="status-bar">
                <div id="cameras"></div>
                <span id="status"></span>
            </div>`;

        this.countdownEl = document.getElementById("countdown");
        this.statusBtnEl = document.getElementById("status-button");
        this.statusBtnErrorEl = document.getElementById("status-button-error");
        this.camerasDivEl = document.getElementById("cameras");
        this.statusEl = document.getElementById("status");

        this.spinnerEl.style.display = "none";
        this.spinnerEl.style.marginBottom = "";

        this.statusBtnEl.addEventListener("click", this.handleStatusButton.bind(this));
        this.statusEl.addEventListener("click", this.handleXmark.bind(this));
    }

    error(message, spinner) {
        this.bodyEl.innerHTML = `
            <fieldset style="margin: 0 30px">
                <p class="error">${message}</p>
            </fieldset>`;

        this.spinnerEl.style.display = spinner ? "inline-block" : "none";
        this.spinnerEl.style.marginBottom = spinner ? "20px" : "";
    }

    async update() {
        const states = await get("data");
        console.log("states:", states);
        const { name, heartbeat, webpoint, timeLeft, cameras, error } = states;
        if (error) throw new Error(error);

        if (!this.isHealthy) this.render(name);
        this.isHealthy = true;

        this.updateWebpoint(webpoint, heartbeat);
        this.updateCameras(cameras);

        this.updateCountdown(timeLeft);
        this.updateStatusButton(timeLeft);
    }

    handleError(err) {
        console.error(`handle error: ${err.message}`);
        let message = err.message;
        let spinner = false;
        if (err instanceof TypeError) {
            message = "Нет соединения с веб-сервером";
            spinner = true;
        }

        this.isHealthy = false
        this.error(message, spinner)
    }

    updateWebpoint(webpoint, heartbeat) {
        if (this.statusEl.className === circleXmarkClassName) {
            this.statusEl.style.display = "inline-block";
            return;
        }

        if (!webpoint) {
            this.statusEl.className = linkSlashClassName;
            this.statusEl.style.display = "inline-block";
            return
        }

        if (!heartbeat) {
            this.statusEl.className = heartCrackClassName;
            this.statusEl.style.display = "inline-block";
            return
        }

        this.statusEl.className = "";
        this.statusEl.style.display = "none";
    }

    updateCountdown(timeLeft = 0) {
        const hours   = App.formatNumber(Math.floor(timeLeft / 3600));
        const minutes = App.formatNumber(Math.floor(timeLeft / 60 - hours * 60));
        const seconds = App.formatNumber(timeLeft % 60);

        this.countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
        this.setStatus(this.countdownEl, !timeLeft);
    }

    static formatNumber = (num) => num < 10 ? "0" + num: num;

    updateStatusButton(timeLeft = 0) {
        this.setStatus(this.statusBtnEl, !timeLeft);
        this.statusBtnEl.disabled = !!timeLeft;
    }

    async handleStatusButton() {
        try {
            await get("button-press", "text");
            this.statusEl.className = circleCheckClassName;
            this.statusEl.style.display = "inline-block";
        } catch (err) {
            this.statusEl.className = circleXmarkClassName;
            this.statusEl.style.display = "inline-block";
            this.statusBtnEl.style.display = "none";
            this.statusBtnErrorEl.style.display = "flex";
            this.statusBtnErrorEl.firstElementChild.innerText = err.message;
        }
    }

    handleXmark() {
        if (this.statusEl.className !== circleXmarkClassName) {
            return
        }

        this.statusBtnEl.style.display = "block";
        this.statusBtnErrorEl.style.display = "none";
        this.statusEl.className = "";
        this.statusEl.style.display = "none";
    }

    createCamera(cameraID, states) {
        console.log(`creating camera ${cameraID}...`);
        const {name, car, human, inputs} = states;
        const camera = newElement("fieldset", { className: "camera", id: cameraID });
        const legend = newElement("legend", { innerText: name });

        const truckIcon = newElement("span", { className: truckIconClassName });
        const humanIcon = newElement("span", { className: humanIconClassName });

        let inputNum = 0;
        const inputIcons = Object.entries(inputs).map(([name, inp]) => {
            const inputEl = newElement("span", {
                className: inputIconClassName + " tooltip",
                id: inp.id
            })
            inputEl.appendChild(newElement("span", {
                className: "tooltip-text",
                // innerText: `Input: ${name}`,
                innerText: `Input: ${++inputNum}`,
            }))

            return inputEl;
        });

        this.setStatus(truckIcon, car);
        this.setStatus(humanIcon, human);
        inputIcons.forEach(icon => this.setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state));

        [legend, truckIcon, humanIcon, ...inputIcons].forEach(el => camera.appendChild(el));
        this.camerasDivEl.appendChild(camera);
        return camera;
    }

    updateCameras(cameras) {
        const prevCameras = this.cameras;
        this.cameras = cameras;

        const toRemove = Object.keys(prevCameras).filter(prevID => !Object.keys(this.cameras).find(currID => prevID === currID));
        toRemove.forEach(id => {
            console.log(`removing camera ${id} states...`);
            const el = document.getElementById(id);
            el?.parentNode.removeChild(el);
        });

        Object.entries(cameras).forEach(([id, states]) => {
            const camera = document.getElementById(id);
            camera ? this.updateCamera(camera, states) : this.createCamera(id, states);
        })
    }

    updateCamera(camera, states) {
        console.log(`updating camera ${camera.id} states...`);
        const {car, human, inputs} = states;
        for (const icon of camera.getElementsByTagName("span")) { // todo: pay attention
            if (icon.className.includes(truckIconClassName)) {
                this.setStatus(icon, car);
            } else if (icon.className.includes(humanIconClassName)) {
                this.setStatus(icon, human);
            } else if (icon.className.includes(inputIconClassName)) {
                this.setStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state);
            }
        }
    }

    setStatus(element, status) {
        const isReady = element.classList.contains("ready");
        if (isReady && !status || !isReady && status) {
            element.classList.toggle("ready");
        }
    }
}

window.onload = () => {
    const app = new App();
    setTimeout(async function cycle() {
        try {
            await app.update();
        } catch (err) {
            app.handleError(err);
        }

        setTimeout(cycle, pollingFrequency);
    });
};
