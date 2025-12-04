// Smooth scrolling for navigation links
document.querySelectorAll('.nav-links a').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const targetId = this.getAttribute('href');
        const targetSection = document.querySelector(targetId);

        if (targetSection) {
            const navbarHeight = document.querySelector('.navbar').offsetHeight;
            const targetPosition = targetSection.offsetTop - navbarHeight;

            window.scrollTo({
                top: targetPosition,
                behavior: 'smooth'
            });
        }
    });
});

// Navbar scroll effect
let lastScroll = 0;
const navbar = document.querySelector('.navbar');

window.addEventListener('scroll', () => {
    const currentScroll = window.pageYOffset;

    if (currentScroll > 100) {
        navbar.style.boxShadow = '0 4px 20px rgba(0, 0, 0, 0.3)';
    } else {
        navbar.style.boxShadow = 'none';
    }

    lastScroll = currentScroll;
});

// Intersection Observer for fade-in animations
const observerOptions = {
    threshold: 0.1,
    rootMargin: '0px 0px -100px 0px'
};

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Observe feature cards
document.querySelectorAll('.feature-card').forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(30px)';
    card.style.transition = `all 0.6s ease ${index * 0.1}s`;
    observer.observe(card);
});

// Observe tech cards
document.querySelectorAll('.tech-card').forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(30px)';
    card.style.transition = `all 0.6s ease ${index * 0.1}s`;
    observer.observe(card);
});

// Observe step cards
document.querySelectorAll('.step-card').forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(30px)';
    card.style.transition = `all 0.6s ease ${index * 0.1}s`;
    observer.observe(card);
});

// Observe layer cards
document.querySelectorAll('.layer-card').forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateX(-30px)';
    card.style.transition = `all 0.6s ease ${index * 0.15}s`;
    observer.observe(card);
});

// Button click animations
document.querySelectorAll('.btn').forEach(button => {
    button.addEventListener('click', function(e) {
        // Create ripple effect
        const ripple = document.createElement('span');
        const rect = this.getBoundingClientRect();
        const size = Math.max(rect.width, rect.height);
        const x = e.clientX - rect.left - size / 2;
        const y = e.clientY - rect.top - size / 2;

        ripple.style.width = ripple.style.height = size + 'px';
        ripple.style.left = x + 'px';
        ripple.style.top = y + 'px';
        ripple.classList.add('ripple');

        this.appendChild(ripple);

        setTimeout(() => {
            ripple.remove();
        }, 600);
    });
});

// Add ripple styles dynamically
const style = document.createElement('style');
style.textContent = `
    .btn {
        position: relative;
        overflow: hidden;
    }

    .ripple {
        position: absolute;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.3);
        transform: scale(0);
        animation: ripple-animation 0.6s ease-out;
        pointer-events: none;
    }

    @keyframes ripple-animation {
        to {
            transform: scale(4);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// Code snippet copy functionality
document.querySelectorAll('.code-snippet').forEach(snippet => {
    const copyButton = document.createElement('button');
    copyButton.innerHTML = '📋';
    copyButton.className = 'copy-button';
    copyButton.title = 'Copy to clipboard';

    copyButton.addEventListener('click', async () => {
        const code = snippet.querySelector('code').textContent;
        try {
            await navigator.clipboard.writeText(code);
            copyButton.innerHTML = '✅';
            copyButton.style.background = 'var(--success-color)';

            setTimeout(() => {
                copyButton.innerHTML = '📋';
                copyButton.style.background = '';
            }, 2000);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
    });

    snippet.style.position = 'relative';
    snippet.appendChild(copyButton);
});

// Add copy button styles
const copyButtonStyle = document.createElement('style');
copyButtonStyle.textContent = `
    .copy-button {
        position: absolute;
        top: 8px;
        right: 8px;
        background: rgba(0, 173, 216, 0.2);
        border: 1px solid rgba(0, 173, 216, 0.3);
        color: var(--primary-color);
        padding: 6px 12px;
        border-radius: 6px;
        cursor: pointer;
        font-size: 14px;
        transition: all 0.3s ease;
        z-index: 10;
    }

    .copy-button:hover {
        background: rgba(0, 173, 216, 0.3);
        transform: scale(1.05);
    }

    .copy-button:active {
        transform: scale(0.95);
    }
`;
document.head.appendChild(copyButtonStyle);

// Parallax effect for hero section
window.addEventListener('scroll', () => {
    const scrolled = window.pageYOffset;
    const heroVisual = document.querySelector('.hero-visual');

    if (heroVisual) {
        heroVisual.style.transform = `translateY(${scrolled * 0.3}px)`;
    }
});

// Add typing effect to hero title
const heroTitle = document.querySelector('.hero-title');
if (heroTitle) {
    const originalHTML = heroTitle.innerHTML;
    heroTitle.style.opacity = '1';

    // Optional: Add cursor blink effect
    const gradientText = heroTitle.querySelector('.gradient-text');
    if (gradientText) {
        let isVisible = true;
        setInterval(() => {
            if (isVisible) {
                gradientText.style.opacity = '0.7';
            } else {
                gradientText.style.opacity = '1';
            }
            isVisible = !isVisible;
        }, 2000);
    }
}

// Add hover effect to code window
const codeWindows = document.querySelectorAll('.code-window');
codeWindows.forEach(window => {
    window.addEventListener('mouseenter', function() {
        this.style.transform = 'scale(1.02)';
        this.style.boxShadow = '0 12px 40px rgba(0, 173, 216, 0.3)';
    });

    window.addEventListener('mouseleave', function() {
        this.style.transform = 'scale(1)';
        this.style.boxShadow = 'var(--shadow-lg)';
    });
});

// Add transition to code windows
codeWindows.forEach(window => {
    window.style.transition = 'all 0.3s ease';
});

// Easter egg: Konami code
let konamiCode = [];
const konamiSequence = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a'];

document.addEventListener('keydown', (e) => {
    konamiCode.push(e.key);
    konamiCode = konamiCode.slice(-10);

    if (konamiCode.join('') === konamiSequence.join('')) {
        // Activate easter egg
        document.body.style.animation = 'rainbow 2s linear infinite';

        const easterEggStyle = document.createElement('style');
        easterEggStyle.textContent = `
            @keyframes rainbow {
                0% { filter: hue-rotate(0deg); }
                100% { filter: hue-rotate(360deg); }
            }
        `;
        document.head.appendChild(easterEggStyle);

        setTimeout(() => {
            document.body.style.animation = '';
        }, 5000);
    }
});

// Performance optimization: Lazy load images
if ('IntersectionObserver' in window) {
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.classList.remove('lazy');
                imageObserver.unobserve(img);
            }
        });
    });

    document.querySelectorAll('img[data-src]').forEach(img => {
        imageObserver.observe(img);
    });
}

// Add active state to navigation links based on scroll position
window.addEventListener('scroll', () => {
    const sections = document.querySelectorAll('section[id]');
    const navLinks = document.querySelectorAll('.nav-links a');

    let current = '';

    sections.forEach(section => {
        const sectionTop = section.offsetTop;
        const sectionHeight = section.clientHeight;

        if (window.pageYOffset >= sectionTop - 200) {
            current = section.getAttribute('id');
        }
    });

    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('href') === `#${current}`) {
            link.classList.add('active');
        }
    });
});

// Add active link style
const activeLinkStyle = document.createElement('style');
activeLinkStyle.textContent = `
    .nav-links a.active {
        color: var(--primary-color);
    }

    .nav-links a.active::after {
        width: 100%;
    }
`;
document.head.appendChild(activeLinkStyle);

console.log('🚀 Go API Skeleton - Website loaded successfully!');
console.log('Built with ❤️ using vanilla JavaScript');
