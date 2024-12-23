let mediaRecorder;
let audioChunks = [];

document.getElementById('startRecording').addEventListener('click', async () => {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    mediaRecorder = new MediaRecorder(stream);

    mediaRecorder.start();
    document.getElementById('stopRecording').disabled = false;
    document.getElementById('startRecording').disabled = true;

    mediaRecorder.ondataavailable = (event) => {
        audioChunks.push(event.data);
    };

    mediaRecorder.onstop = async () => {
        const audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
        const audioUrl = URL.createObjectURL(audioBlob);
        document.getElementById('audioPlayback').src = audioUrl;

        // Send the audio to the backend
        const formData = new FormData();
        formData.append('audio', audioBlob, 'recording.wav');

        await fetch('http://localhost:8080/upload', {
            method: 'POST',
            body: formData
        });

        // Reset for next recording
        audioChunks = [];
        document.getElementById('startRecording').disabled = false;
        document.getElementById('stopRecording').disabled = true;
    };
});

document.getElementById('stopRecording').addEventListener('click', () => {
    mediaRecorder.stop();
});
