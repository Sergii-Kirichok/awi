const green = "#13a017";
const red   = "#ff0000";

const h = 0;
const m = 0;
const s = 5;

let time;

const countdownEl = document.getElementById("countdown");
const statusBtnEl = document.getElementById("status-button");

statusBtnEl.addEventListener("click", () => { if (!time) recover() })

function recover() {
    statusBtnEl.disabled = true;
    statusBtnEl.style.backgroundColor = red;
    time = (h * 60 + m) * 60 + s;
    countdown();
}

recover();

function countdown() {
    setTimeout(function again() {
        updateCountdown()
        if (!time) {
            statusBtnEl.style.backgroundColor = green;
            statusBtnEl.innerText = "Взвешивание разрешено";
            statusBtnEl.disabled = false;
            return
        }

        setTimeout(again, 1000)
    }, 1000)
}

function updateCountdown() {
    time--;

    const hours = decrease(Math.floor(time / 3600));
    const minutes = decrease(Math.floor(time / 60 - hours * 60));
    const seconds = decrease(time % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

decrease = (num) => num < 10 ? "0" + num: num