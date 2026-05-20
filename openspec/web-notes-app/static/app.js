async function submitNote() {
    const title = document.getElementById('noteTitle').value.trim();
    const content = document.getElementById('noteContent').value.trim();
    const tags = document.getElementById('noteTags').value.trim();
    const errorDiv = document.getElementById('createError');

    errorDiv.classList.add('d-none');
    errorDiv.textContent = '';

    if (!title) {
        errorDiv.textContent = 'Title is required';
        errorDiv.classList.remove('d-none');
        return;
    }

    if (!content) {
        errorDiv.textContent = 'Content is required';
        errorDiv.classList.remove('d-none');
        return;
    }

    try {
        const formData = new FormData();
        formData.append('title', title);
        formData.append('content', content);
        formData.append('tags', tags);

        const response = await fetch('/note', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Failed to create note');
        }

        const html = await response.text();
        document.getElementById('notesBoard').innerHTML = html;

        const modal = bootstrap.Modal.getInstance(document.getElementById('createModal'));
        modal.hide();

        document.getElementById('createForm').reset();
    } catch (error) {
        errorDiv.textContent = error.message;
        errorDiv.classList.remove('d-none');
    }
}

async function viewNote(id) {
    try {
        const response = await fetch('/note/' + id);
        if (!response.ok) {
            throw new Error('Failed to load note');
        }

        const html = await response.text();
        const modal = new bootstrap.Modal(document.getElementById('viewModal'));
        document.getElementById('viewModal').querySelector('.modal-content').innerHTML = html;
        modal.show();
    } catch (error) {
        alert('Failed to load note: ' + error.message);
    }
}

document.getElementById('createForm').addEventListener('keypress', function(e) {
    if (e.key === 'Enter' && e.target.tagName !== 'TEXTAREA') {
        e.preventDefault();
        submitNote();
    }
});