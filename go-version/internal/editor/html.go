package editor

const editorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>qkflow - Add PR Description</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/easymde@2.18.0/dist/easymde.min.css">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            background: #0d1117;
            color: #c9d1d9;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 6px;
            overflow: hidden;
        }
        
        .header {
            padding: 16px 24px;
            background: #0d1117;
            border-bottom: 1px solid #30363d;
        }
        
        .header h1 {
            font-size: 20px;
            font-weight: 600;
            color: #c9d1d9;
        }
        
        .header p {
            font-size: 14px;
            color: #8b949e;
            margin-top: 4px;
        }
        
        .editor-container {
            padding: 24px;
        }
        
        .EasyMDEContainer .CodeMirror {
            background: #0d1117;
            color: #c9d1d9;
            border: 1px solid #30363d;
            border-radius: 6px;
            min-height: 300px;
        }
        
        .EasyMDEContainer .CodeMirror-scroll {
            min-height: 300px;
        }
        
        .editor-toolbar {
            background: #161b22;
            border: 1px solid #30363d;
            border-bottom: none;
            border-radius: 6px 6px 0 0;
        }
        
        .editor-toolbar button {
            color: #c9d1d9 !important;
        }
        
        .editor-toolbar button:hover {
            background: #30363d;
            border-color: #8b949e;
        }
        
        .editor-toolbar.fullscreen {
            background: #0d1117;
        }
        
        .upload-section {
            margin-top: 20px;
            padding: 20px;
            background: #0d1117;
            border: 2px dashed #30363d;
            border-radius: 6px;
            text-align: center;
            transition: all 0.3s;
        }
        
        .upload-section.dragover {
            border-color: #58a6ff;
            background: #161b22;
        }
        
        .upload-section h3 {
            font-size: 16px;
            margin-bottom: 8px;
            color: #c9d1d9;
        }
        
        .upload-section p {
            font-size: 14px;
            color: #8b949e;
            margin-bottom: 12px;
        }
        
        .upload-btn {
            padding: 8px 16px;
            background: #21262d;
            border: 1px solid #30363d;
            border-radius: 6px;
            color: #c9d1d9;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.2s;
        }
        
        .upload-btn:hover {
            background: #30363d;
        }
        
        .file-list {
            margin-top: 16px;
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        
        .file-item {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 12px;
            background: #0d1117;
            border: 1px solid #30363d;
            border-radius: 6px;
        }
        
        .file-info {
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .file-icon {
            font-size: 24px;
        }
        
        .file-details {
            display: flex;
            flex-direction: column;
        }
        
        .file-name {
            font-size: 14px;
            color: #c9d1d9;
            font-weight: 500;
        }
        
        .file-size {
            font-size: 12px;
            color: #8b949e;
        }
        
        .remove-btn {
            padding: 4px 12px;
            background: #21262d;
            border: 1px solid #30363d;
            border-radius: 6px;
            color: #f85149;
            cursor: pointer;
            font-size: 12px;
            transition: all 0.2s;
        }
        
        .remove-btn:hover {
            background: #f85149;
            color: #fff;
        }
        
        .actions {
            padding: 16px 24px;
            background: #0d1117;
            border-top: 1px solid #30363d;
            display: flex;
            justify-content: flex-end;
            gap: 12px;
        }
        
        .btn {
            padding: 8px 20px;
            border: none;
            border-radius: 6px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .btn-cancel {
            background: #21262d;
            color: #c9d1d9;
            border: 1px solid #30363d;
        }
        
        .btn-cancel:hover {
            background: #30363d;
        }
        
        .btn-save {
            background: #238636;
            color: #fff;
        }
        
        .btn-save:hover {
            background: #2ea043;
        }
        
        .btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
        
        .hint {
            margin-top: 16px;
            padding: 12px;
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 6px;
            font-size: 13px;
            color: #8b949e;
        }
        
        .hint strong {
            color: #c9d1d9;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📝 Add PR Description & Screenshots</h1>
            <p>Write your description in Markdown and drag & drop images or videos</p>
        </div>
        
        <div class="editor-container">
            <textarea id="editor"></textarea>
            
            <div class="upload-section" id="uploadSection">
                <h3>📎 Attach Images & Videos</h3>
                <p>Drag & drop files here, or click to select</p>
                <input type="file" id="fileInput" multiple accept="image/*,video/*" style="display:none">
                <button class="upload-btn" onclick="document.getElementById('fileInput').click()">
                    Choose Files
                </button>
                <div class="file-list" id="fileList"></div>
            </div>
            
            <div class="hint">
                <strong>💡 Tip:</strong> You can paste images directly from clipboard (Cmd+V / Ctrl+V), 
                or drag & drop files from Finder/Explorer. Supported formats: PNG, JPG, GIF, WebP, SVG, MP4, MOV, WebM.
            </div>
        </div>
        
        <div class="actions">
            <button class="btn btn-cancel" onclick="cancel()">Cancel</button>
            <button class="btn btn-save" onclick="save()">Save and Continue</button>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/easymde@2.18.0/dist/easymde.min.js"></script>
    <script>
        // Initialize editor
        const easyMDE = new EasyMDE({
            element: document.getElementById('editor'),
            spellChecker: false,
            placeholder: "## Description\n\nExplain what changes you made and why.\n\n### Screenshots\n\n(Images will be automatically inserted here when you upload them)\n\n### Notes\n\n- Any additional information\n- Dependencies\n- Breaking changes",
            autosave: {
                enabled: true,
                uniqueId: "qkflow-pr-editor",
                delay: 1000,
            },
            toolbar: [
                "bold", "italic", "heading", "|",
                "quote", "unordered-list", "ordered-list", "|",
                "link", "image", "|",
                "preview", "side-by-side", "fullscreen", "|",
                "guide"
            ],
        });

        // Track uploaded files
        let uploadedFiles = [];

        // Drag and drop
        const uploadSection = document.getElementById('uploadSection');
        const fileInput = document.getElementById('fileInput');
        const fileList = document.getElementById('fileList');

        uploadSection.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadSection.classList.add('dragover');
        });

        uploadSection.addEventListener('dragleave', () => {
            uploadSection.classList.remove('dragover');
        });

        uploadSection.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadSection.classList.remove('dragover');
            handleFiles(e.dataTransfer.files);
        });

        fileInput.addEventListener('change', (e) => {
            handleFiles(e.target.files);
        });

        // Paste image from clipboard
        document.addEventListener('paste', (e) => {
            const items = e.clipboardData.items;
            for (let i = 0; i < items.length; i++) {
                if (items[i].type.indexOf('image') !== -1) {
                    const blob = items[i].getAsFile();
                    const file = new File([blob], 'pasted-image-' + Date.now() + '.png', {type: blob.type});
                    handleFiles([file]);
                    e.preventDefault();
                }
            }
        });

        async function handleFiles(files) {
            for (let file of files) {
                await uploadFile(file);
            }
        }

        async function uploadFile(file) {
            const formData = new FormData();
            formData.append('file', file);

            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData
                });

                if (!response.ok) {
                    throw new Error('Upload failed');
                }

                const data = await response.json();
                uploadedFiles.push({
                    name: data.filename,
                    path: data.path,
                    size: data.size
                });

                addFileToList(data.filename, data.size);
                insertImageToEditor(data.filename);
            } catch (error) {
                alert('Failed to upload file: ' + error.message);
            }
        }

        function addFileToList(filename, size) {
            const fileItem = document.createElement('div');
            fileItem.className = 'file-item';
            fileItem.innerHTML = '<div class="file-info">' +
                '<span class="file-icon">' + getFileIcon(filename) + '</span>' +
                '<div class="file-details">' +
                '<span class="file-name">' + filename + '</span>' +
                '<span class="file-size">' + formatSize(size) + '</span>' +
                '</div></div>' +
                '<button class="remove-btn" onclick="removeFile(\'' + filename + '\')">Remove</button>';
            fileList.appendChild(fileItem);
        }

        function insertImageToEditor(filename) {
            const cursor = easyMDE.codemirror.getCursor();
            const isImage = /\.(png|jpg|jpeg|gif|webp|svg)$/i.test(filename);
            const isVideo = /\.(mp4|mov|webm|avi)$/i.test(filename);
            
            let markdown = '';
            if (isImage) {
                markdown = '![Image](./' + filename + ')\n\n';
            } else if (isVideo) {
                markdown = '[Video: ' + filename + '](./' + filename + ')\n\n';
            }
            
            easyMDE.codemirror.replaceRange(markdown, cursor);
        }

        function removeFile(filename) {
            uploadedFiles = uploadedFiles.filter(f => f.name !== filename);
            const fileItems = fileList.children;
            for (let item of fileItems) {
                if (item.textContent.includes(filename)) {
                    item.remove();
                    break;
                }
            }
        }

        function getFileIcon(filename) {
            if (/\.(png|jpg|jpeg|gif|webp|svg)$/i.test(filename)) return '🖼️';
            if (/\.(mp4|mov|webm|avi)$/i.test(filename)) return '🎥';
            return '📄';
        }

        function formatSize(bytes) {
            if (bytes < 1024) return bytes + ' B';
            if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
            return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
        }

        async function save() {
            const content = easyMDE.value();
            const files = uploadedFiles.map(f => f.path);

            try {
                const response = await fetch('/save', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ content, files })
                });

                if (response.ok) {
                    document.body.innerHTML = '<div style="text-align:center;padding:50px;"><h1>✅ Saved!</h1><p>You can close this window now.</p></div>';
                } else {
                    alert('Failed to save');
                }
            } catch (error) {
                alert('Error: ' + error.message);
            }
        }

        async function cancel() {
            if (confirm('Are you sure you want to cancel? All changes will be lost.')) {
                try {
                    await fetch('/cancel', { method: 'POST' });
                    document.body.innerHTML = '<div style="text-align:center;padding:50px;"><h1>❌ Cancelled</h1><p>You can close this window now.</p></div>';
                } catch (error) {
                    window.close();
                }
            }
        }
    </script>
</body>
</html>
`

