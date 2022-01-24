const zone = window.location.pathname.slice("/zones/".length);

const pollingFrequency = 1000; // in ms

const truckIconClassName = "icon truck";
const humanIconClassName = "icon human";
const inputIconClassName = "icon input";

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
            <button id="status-button">Зважити</button>
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
        console.log(states);
        const { name, heartbeat, timeLeft, cameras, error } = states;
        if (error) throw new Error(error);

        if (!this.isHealthy) this.render(name);
        this.isHealthy = true;

        this.updateHeartbeat(heartbeat);
        this.updateCameras(cameras);
        this.updateCountdown(timeLeft);
        this.updateStatusButton(timeLeft);
    }

    handleError(err) {
        console.error(err);
        let message = err.message;
        let spinner = false;
        if (err instanceof TypeError) {
            message = "Нет соединения с веб-сервером";
            spinner = true;
        }

        this.isHealthy = false
        this.error(message, spinner)
    }

    updateHeartbeat(heartbeat) {
        if ([circleCheckClassName, circleXmarkClassName].includes(this.statusEl.className)) return

        if (!heartbeat) {
            this.setStatusIcon(heartCrackClassName);
        } else {
            this.hideStatusIcon();
        }
    }

    setStatusIcon(className) {
        this.statusEl.className = className;
        this.statusEl.style.display = "inline-block";
    }

    hideStatusIcon() {
        this.statusEl.className = "";
        this.statusEl.style.display = "none";
    }

    updateCountdown(timeLeft = 0) {
        const hours   = App.formatNumber(Math.floor(timeLeft / 3600));
        const minutes = App.formatNumber(Math.floor(timeLeft / 60 - hours * 60));
        const seconds = App.formatNumber(timeLeft % 60);

        this.countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
        this.setReadiness(this.countdownEl, !timeLeft);
    }

    static formatNumber = (num) => num < 10 ? "0" + num: num;

    updateStatusButton(timeLeft = 0) {
        this.setReadiness(this.statusBtnEl, !timeLeft);
    }

    async handleStatusButton() {
        try {
            await get("button-press", "text");
        } catch (err) {
            this.setStatusIcon(circleXmarkClassName);
            this.toggleStatusButton(err);
            return
        }

        this.setStatusIcon(circleCheckClassName);
        this.statusBtnEl.disabled = true;
        this.statusBtnEl.classList.add("clicked");
        setTimeout(() => {
            this.hideStatusIcon();
            this.statusBtnEl.disabled = false;
            this.statusBtnEl.classList.remove("clicked");
        }, 10000);
    }

    toggleStatusButton(err) {
        const bool = err != null;
        this.statusBtnEl.style.display = bool ? "none" : "block";
        this.statusBtnErrorEl.style.display = bool ? "flex" : "none";
        if (bool) this.statusBtnErrorEl.firstElementChild.innerText = err.message;
    }

    handleXmark() {
        if (this.statusEl.className !== circleXmarkClassName) return

        this.toggleStatusButton();
        this.hideStatusIcon();
    }

    createCamera(cameraID, states) {
        const {name, car, human, inputs, connection} = states;

        const camera = newElement("div", { className: "camera", id: cameraID });
        const fieldset = newElement("fieldset");
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

        const p = newElement("p", { className: "connection-state", innerText: "Состояние соединения: " });
        const mark = newElement("mark", { className: connection.state ? "green" : "red", innerText: connection.type });

        this.setCameraStatus(truckIcon, car);
        this.setCameraStatus(humanIcon, human);
        inputIcons.forEach(icon => this.setCameraStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state));

        [legend, truckIcon, humanIcon, ...inputIcons].forEach(el => fieldset.appendChild(el));
        p.appendChild(mark);
        camera.appendChild(fieldset);
        camera.appendChild(p);
        this.camerasDivEl.appendChild(camera);
        return camera;
    }

    updateCameras(cameras) {
        const prevCameras = this.cameras;
        this.cameras = cameras;

        const toRemove = Object.keys(prevCameras).filter(prevID => !Object.keys(this.cameras).find(currID => prevID === currID));
        toRemove.forEach(id => {
            const el = document.getElementById(id);
            el?.parentNode.removeChild(el);
        });

        Object.entries(cameras).forEach(([id, states]) => {
            const camera = document.getElementById(id);
            camera ? this.updateCamera(camera, states) : this.createCamera(id, states);
        })
    }

    updateCamera(camera, states) {
        const {car, human, inputs, connection} = states;
        for (const mark of camera.getElementsByTagName("mark")) {
            mark.className = connection.state ? "green" : "red";
            mark.innerText = connection.type;
        }

        for (const icon of camera.getElementsByTagName("span")) {
            if (icon.className.includes(truckIconClassName)) {
                this.setCameraStatus(icon, car);
            } else if (icon.className.includes(humanIconClassName)) {
                this.setCameraStatus(icon, human);
            } else if (icon.className.includes(inputIconClassName)) {
                this.setCameraStatus(icon, Object.values(inputs).find(inp => icon.id === inp.id).state);
            }
        }
    }

    setCameraStatus(element, status) {
        const expander = element.className.split(" ").pop();
        if (["green", "red"].includes(expander)) {
            element.classList.remove(expander);
        }

        if (status) element.classList.add(status);
    }

    setReadiness(element, status) {
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
