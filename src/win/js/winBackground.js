const backgroundSel = document.getElementById('background-select');
const canvas = document.createElement('canvas');
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;
document.body.appendChild(canvas);
const ctx = canvas.getContext('2d');

const particles = [];
const numParticles = 100;

for (let i = 0; i < numParticles; i++) {
    particles.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        radius: Math.random() * 5 + 2,
        color: `hsl(${Math.random() * 360}, 50%, 50%)`,
        speed: Math.random() / 2,
        angle: Math.random() * Math.PI * 2
    });
}

function animateParticles() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    if (backgroundSel.value === "particles") {
        for (let i = 0; i < particles.length; i++) {
            const particle = particles[i];
            particle.x += Math.cos(particle.angle) * particle.speed;
            particle.y += Math.sin(particle.angle) * particle.speed;

            if (particle.x < 0 || particle.x > canvas.width) {
                particle.angle = Math.PI - particle.angle;
            }
            if (particle.y < 0 || particle.y > canvas.height) {
                particle.angle = -particle.angle;
            }

            ctx.beginPath();
            ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2);
            ctx.fillStyle = particle.color;
            ctx.fill();
        }
    }

    requestAnimationFrame(animateParticles);
}

animateParticles();