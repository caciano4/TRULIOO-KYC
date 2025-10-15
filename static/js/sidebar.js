// Enhanced Navigation with Active States
function showPage(pageId) {
    // Remove active class from all pages
    const pages = document.querySelectorAll('.page');
    pages.forEach(page => {
        page.classList.remove('active');
    });

    // Remove active class from all sidebar links
    const navLinks = document.querySelectorAll('.sidebar ul li a');
    navLinks.forEach(link => {
        link.classList.remove('active');
    });

    // Show the selected page
    const activePage = document.getElementById(pageId);
    if (activePage) {
        activePage.classList.add('active');
    }

    // Add active class to the corresponding sidebar link
    const activeLink = document.querySelector(`[onclick="showPage('${pageId}')"]`);
    if (activeLink) {
        activeLink.classList.add('active');
    }

    // Update page title
    const pageTitles = {
        'dashboard': 'Dashboard',
        'uploadKYC': 'Upload KYC',
        'checkPackages': 'Pacotes KYC',
        'checkKYC': 'Clientes KYC'
    };

    if (pageTitles[pageId]) {
        document.title = `${pageTitles[pageId]} - KORE KYC Platform`;
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    // Set default active page
    showPage('dashboard');

    // Add keyboard navigation
    document.addEventListener('keydown', function(e) {
        if (e.altKey) {
            switch(e.key) {
                case '1':
                    e.preventDefault();
                    showPage('dashboard');
                    break;
                case '2':
                    e.preventDefault();
                    showPage('uploadKYC');
                    break;
                case '3':
                    e.preventDefault();
                    showPage('checkPackages');
                    break;
                case '4':
                    e.preventDefault();
                    showPage('checkKYC');
                    break;
            }
        }
    });

    // Add smooth scrolling for page transitions
    const style = document.createElement('style');
    style.textContent = `
        .page {
            transition: opacity 0.3s ease-in-out;
        }
        .page:not(.active) {
            opacity: 0;
            pointer-events: none;
        }
        .page.active {
            opacity: 1;
            pointer-events: all;
        }
    `;
    document.head.appendChild(style);
});

