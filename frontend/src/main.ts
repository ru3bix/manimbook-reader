import "./style.css";
import "./app.css";

import { EventsOn } from "../wailsjs/runtime/runtime";
import { GetBook, GetPort } from "../wailsjs/go/main/App";
import { book } from "../wailsjs/go/models";

let currentChapter: number | null = 0;
let port: number | null = null;

EventsOn("bookOpen", () => {
  Promise.all([GetBook(), GetPort()]).then((res) => {
    const book = res[0];
    port = res[1];
    currentChapter = 0;
    removeChapters();
    appendChapters(book.chapters, book);
    changeChapter(currentChapter, book);
  });
});

function changeChapter(ch: number, book: book.Book) {
  currentChapter = ch;
  const title = document.getElementById("chapterTitle");
  const frame = document.getElementById("chapter") as HTMLIFrameElement;
  const contentsBtn = document.getElementById(`btn-chapter-${ch}`);
  if (title) {
    title.innerText = book.chapters[ch];
  }
  if (contentsBtn) {
    const prevActive = document.querySelector(".active");
    prevActive?.classList.remove("active");
    contentsBtn.classList.add("active");
  }
  frame.src = `http://127.0.0.1:${port}/book/${book.chapters[ch]}/index.html`;
}

EventsOn("incZoom", (data: number) => {
  ZoomIn(data);
});
EventsOn("decZoom", (data: number) => {
  ZoomOut(data);
});
EventsOn("setZoom", (data: number) => {
  setZoom(data);
});

EventsOn("showContents", showContents);
EventsOn("hideContents", hideContents);
EventsOn("showNavbar", showNavbar);
EventsOn("hideNavbar", hideNavbar);

function appendChapters(chapters: string[], book: book.Book) {
  const contentsDiv = document.getElementById("contents");
  if (!contentsDiv) return;

  chapters.forEach((chapter, index) => {
    const btn = document.createElement("button");
    btn.id = `btn-chapter-${index}`;
    btn.onclick = () => {
      changeChapter(index, book);
    };
    btn.textContent = chapter;
    contentsDiv.appendChild(btn);
  });
}

function removeChapters() {
  const contents = document.getElementById("contents");
  if (contents) {
    contents.querySelectorAll("button").forEach((button) => button.remove());
  }
}

const nxtBtn = document.getElementById("nxtBtn");
const prevBtn = document.getElementById("prevBtn");

if (nxtBtn && prevBtn) {
  nxtBtn.onclick = async () => {
    const book = await GetBook();
    if (currentChapter !== null && book.chapters.length > currentChapter) {
      currentChapter += 1;
      changeChapter(currentChapter, book);
    }
  };
  prevBtn.onclick = async () => {
    const book = await GetBook();
    if (currentChapter && book.chapters.length !== 0) {
      currentChapter -= 1;
      changeChapter(currentChapter, book);
    }
  };
}

const openBtn = document.getElementById("openContents");
const closeBtn = document.getElementById("closeContents");

if (openBtn && closeBtn) {
  openBtn.onclick = showContents;
  closeBtn.onclick = hideContents;
}

function ZoomIn(n: number) {
  const r = document.querySelector(":root") as HTMLElement;
  if (r) {
    const rs = getComputedStyle(r);
    const zoom = rs.getPropertyValue("--scale-factor") ?? 1.0;
    r.style.setProperty("--scale-factor", (parseFloat(zoom) + n).toString());
  }
}

function ZoomOut(n: number) {
  const r = document.querySelector(":root") as HTMLElement;
  if (r) {
    const rs = getComputedStyle(r);
    const zoom = rs.getPropertyValue("--scale-factor") ?? 1.0;
    r.style.setProperty("--scale-factor", (parseFloat(zoom) - n).toString());
  }
}

function setZoom(n: number) {
  const root = document.querySelector(":root") as HTMLElement;
  root.style.setProperty("--scale-factor", n.toString());
}

function showNavbar() {
  const navbar = document.querySelector("nav");
  const container = document.querySelector(".container");
  if (navbar && container) {
    navbar.classList.remove("hidden");
    container.classList.remove("maximized");
  }
}

function hideNavbar() {
  const navbar = document.querySelector("nav");
  const container = document.querySelector(".container");
  if (navbar && container) {
    navbar.classList.add("hidden");
    container.classList.add("maximized");
  }
}

function showContents() {
  const contents = document.getElementById("contents");
  const app = document.getElementById("app");
  if (openBtn && contents && app) {
    openBtn.classList.add("hidden");
    contents.style.width = "250px";
    app.style.marginLeft = "250px";
  }
}

function hideContents() {
  const contents = document.getElementById("contents");
  const app = document.getElementById("app");
  if (openBtn && contents && app) {
    openBtn.classList.remove("hidden");
    contents.style.width = "0px";
    app.style.marginLeft = "0px";
  }
}
