import { rerenderTable } from "./datatable";

function showPage(pageId) {
    console.log(pageId)
    rerenderTable()
    const pages = document.querySelectorAll('.page');
    pages.forEach(page => {
        page.classList.remove('active');
    });

    const activePage = document.getElementById(pageId);
    if (activePage) {
        activePage.classList.add('active');
    }

}

