document.addEventListener('DOMContentLoaded', () => {
    const carouselContainers = document.querySelectorAll('.carousel-container');

    carouselContainers.forEach(container => {
        const slides = container.querySelector('.carousel-slide');
        const images = container.querySelectorAll('.carousel-slide img');
        const prevButton = container.querySelector('.carousel-button.prev');
        const nextButton = container.querySelector('.carousel-button.next');

        if (images.length === 0) {
            console.error("Aucune image trouvée dans le carrousel.", container);
            return;
        }

        let currentIndex = 0;
        let imageWidth = 0;

        // Function to update carousel position
        function updateCarousel() {
            slides.style.transform = `translateX(${-currentIndex * imageWidth}px)`;
        }

        // Set initial image width and start carousel
        const initializeCarousel = () => {
            // Ensure imageWidth is calculated based on the first image of THIS carousel
            imageWidth = images[0].clientWidth;
            updateCarousel();

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
            }, 2500); // Change image every 2.5 seconds
        };

        // Wait for the first image of the current carousel to load
        images[0].onload = initializeCarousel;

        // If images are already cached or loaded, manually trigger initialization
        if (images[0].complete) {
            initializeCarousel();
        }
    });
}); 