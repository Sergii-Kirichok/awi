const countdownEl = document.getElementById("countdown");
const statusBtnEl = document.getElementById("status-button");

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
        await post("/post/reset-timer");
        countdown();
    } catch (error) {
        console.error("can't reset timer");
    } finally {
        disableStatusButton();
    }
}

window.onload = countdown;
statusBtnEl.onclick = async () => await recover();

function countdown() {
    setTimeout(async function again() {
        try {
            const timeLeft = await get("https://sanya.avigilon/get/countdown");
            updateCountdown(timeLeft);
            if (!timeLeft) {
                enableStatusButton();
                return
            }
        } catch (error) {
            console.error("countdown error:", error)
            // todo: mb make something...
            return
        }

        setTimeout(again, 1000);
    });
}

function updateCountdown(timeLeft = 0) {
    const hours = format(Math.floor(timeLeft / 3600));
    const minutes = format(Math.floor(timeLeft / 60 - hours * 60));
    const seconds = format(timeLeft % 60);

    countdownEl.innerText = `${hours}:${minutes}:${seconds}`;
}

format = (num) => num < 10 ? "0" + num: num