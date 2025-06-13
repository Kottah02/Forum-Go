document.addEventListener('DOMContentLoaded', () => {
    const slides = document.querySelector('.carousel-slide');
    const images = document.querySelectorAll('.carousel-slide img');
    const prevButton = document.querySelector('.carousel-button.prev');
    const nextButton = document.querySelector('.carousel-button.next');

    if (images.length === 0) {
        console.error("Aucune image trouvée dans le carrousel.");
        return;
    }

    // Attendre que la première image soit chargée pour obtenir la bonne largeur
    images[0].onload = () => {
        let currentIndex = 0;
        const imageWidth = images[0].clientWidth;

        function updateCarousel() {
            slides.style.transform = `translateX(${-currentIndex * imageWidth}px)`;
        }

        prevButton.addEventListener('click', () => {
            currentIndex = (currentIndex > 0) ? currentIndex - 1 : images.length - 1;
            updateCarousel();
        });

        nextButton.addEventListener('click', () => {
            currentIndex = (currentIndex < images.length - 1) ? currentIndex + 1 : 0;
            updateCarousel();
        });

        // Optional: Auto-play carousel
        setInterval(() => {
            currentIndex = (currentIndex < images.length - 1) ? currentIndex + 1 : 0;
            updateCarousel();
        }, 2500); // Change image every 5 seconds
    };

    // Si les images sont déjà en cache ou chargées, déclencher manuellement l'onload
    if (images[0].complete) {
        images[0].onload();
    }
}); 