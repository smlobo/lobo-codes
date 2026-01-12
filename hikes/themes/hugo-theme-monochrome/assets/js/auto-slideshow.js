
document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll(".auto-slideshow").forEach(slideshow => {
        const slides = slideshow.querySelectorAll(".slide");
        const dots = slideshow.querySelectorAll(".dot");
        const interval = parseInt(slideshow.dataset.interval, 10) || 3000;

        let index = 0;
        let timer;

        function showSlide(i) {
            slides.forEach((slide, idx) => {
                slide.classList.toggle("active", idx === i);
            });

            dots.forEach((dot, idx) => {
                dot.classList.toggle("active", idx === i);
            });

            index = i;
        }

        function nextSlide() {
            showSlide((index + 1) % slides.length);
        }

        function start() {
            timer = setInterval(nextSlide, interval);
        }

        function stop() {
            clearInterval(timer);
        }

        // Dot navigation
        dots.forEach(dot => {
            dot.addEventListener("click", () => {
                stop();
                showSlide(parseInt(dot.dataset.index, 10));
                start();
            });
        });

        // Auto start
        start();

        // Pause on hover
        slideshow.addEventListener("mouseenter", stop);
        slideshow.addEventListener("mouseleave", start);

        // Touch support (NEW)
        slideshow.addEventListener("touchstart", () => {
            isPausedByTouch = true;
            stop();
        }, { passive: true });

        slideshow.addEventListener("touchend", () => {
            if (isPausedByTouch) {
                isPausedByTouch = false;
                start();
            }
        });

        slideshow.addEventListener("touchcancel", () => {
            if (isPausedByTouch) {
                isPausedByTouch = false;
                start();
            }
        });

    });
});
