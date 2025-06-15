
// Animations pour améliorer l'expérience utilisateur

// Animation de fade-in pour les éléments au chargement de la page
document.addEventListener('DOMContentLoaded', () => {
    // Animation pour les posts
    const posts = document.querySelectorAll('.post');
    posts.forEach((post, index) => {
        post.style.opacity = '0';
        post.style.transform = 'translateY(20px)';
        setTimeout(() => {
            post.style.transition = 'all 0.5s ease-out';
            post.style.opacity = '1';
            post.style.transform = 'translateY(0)';
        }, index * 100);
    });

    // Animation pour les formulaires
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.style.opacity = '0';
        form.style.transform = 'scale(0.95)';
        setTimeout(() => {
            form.style.transition = 'all 0.3s ease-out';
            form.style.opacity = '1';
            form.style.transform = 'scale(1)';
        }, 200);
    });

    // Animation pour les boutons
    const buttons = document.querySelectorAll('button, .btn');
    buttons.forEach(button => {
        button.addEventListener('mouseenter', () => {
            button.style.transform = 'scale(1.05)';
            button.style.transition = 'transform 0.2s ease';
        });
        button.addEventListener('mouseleave', () => {
            button.style.transform = 'scale(1)';
        });
    });

    // Animation pour les messages flash
    const flashMessages = document.querySelectorAll('.flash-message');
    flashMessages.forEach(message => {
        message.style.opacity = '0';
        message.style.transform = 'translateY(-20px)';
        setTimeout(() => {
            message.style.transition = 'all 0.3s ease-out';
            message.style.opacity = '1';
            message.style.transform = 'translateY(0)';
        }, 100);

        // Auto-hide après 5 secondes
        setTimeout(() => {
            message.style.opacity = '0';
            message.style.transform = 'translateY(-20px)';
            setTimeout(() => message.remove(), 300);
        }, 5000);
    });
});

// Animation pour les commentaires
function animateNewComment(commentElement) {
    commentElement.style.opacity = '0';
    commentElement.style.transform = 'translateX(-20px)';
    setTimeout(() => {
        commentElement.style.transition = 'all 0.3s ease-out';
        commentElement.style.opacity = '1';
        commentElement.style.transform = 'translateX(0)';
    }, 100);
}

// Animation pour la navigation
document.querySelectorAll('nav a').forEach(link => {
    link.addEventListener('click', (e) => {
        if (!link.classList.contains('active')) {
            e.preventDefault();
            document.body.style.opacity = '0';
            document.body.style.transition = 'opacity 0.3s ease-out';
            
            setTimeout(() => {
                window.location = link.href;
            }, 300);
        }
    });
});

// Animation pour les likes
document.querySelectorAll('.like-button').forEach(button => {
    button.addEventListener('click', () => {
        button.style.transform = 'scale(1.2)';
        setTimeout(() => {
            button.style.transform = 'scale(1)';
        }, 200);
    });
});

// Animation pour les champs de formulaire
document.querySelectorAll('input, textarea').forEach(field => {
    field.addEventListener('focus', () => {
        field.style.transform = 'scale(1.02)';
        field.style.transition = 'transform 0.2s ease';
    });
    
    field.addEventListener('blur', () => {
        field.style.transform = 'scale(1)';
    });
}); 