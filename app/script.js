const copyStatus = document.querySelector('#copy-status');

document.querySelectorAll('[data-copy-target]').forEach((button) => {
    button.addEventListener('click', async () => {
        const source = document.getElementById(button.dataset.copyTarget);
        if (!source) return;

        try {
            await navigator.clipboard.writeText(source.textContent.trim());
            button.textContent = 'Copied';
            button.dataset.copied = 'true';
            copyStatus.textContent = 'Command copied to clipboard.';

            window.setTimeout(() => {
                button.textContent = 'Copy';
                delete button.dataset.copied;
            }, 1800);
        } catch {
            button.textContent = 'Select text';
            copyStatus.textContent = 'Clipboard unavailable. Select and copy the command text.';

            window.setTimeout(() => {
                button.textContent = 'Copy';
            }, 2400);
        }
    });
});
