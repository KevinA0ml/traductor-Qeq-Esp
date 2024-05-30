document.addEventListener('DOMContentLoaded', () => {
    // Sticky header on scroll
    window.addEventListener('scroll', () => {
        const header = document.querySelector('header');
        header.classList.toggle('sticky', window.scrollY > 0);
    });

    // Toggle translation direction
    document.getElementById('cambioidioma').addEventListener('click', toggleLanguage);
    document.getElementById('cambiofrases').addEventListener('click', togglePhrases);
    document.getElementById('clearButton').addEventListener('click', clearInputs);
    document.getElementById('inputText1').addEventListener('input', fetchSuggestions);
    document.getElementById('translateButton').addEventListener('click', translateText);
    document.getElementById('wordBankButton').addEventListener('click', fetchWordBank);

    function toggleLanguage() {
        const inputText1 = document.getElementById('inputText1');
        const inputText2 = document.getElementById('inputText2');
        
        const isQeqchi = inputText1.placeholder === 'Qeqchi';
        inputText1.disabled = false;
        inputText2.disabled = true;
        inputText1.placeholder = isQeqchi ? 'Español' : 'Qeqchi';
        inputText2.placeholder = isQeqchi ? 'Qeqchi' : 'Español';
        clearInputs();
    }

    function togglePhrases() {
        const inputText1 = document.getElementById('inputText1');
        const inputText2 = document.getElementById('inputText2');
        
        const isQeqchiPhrase = inputText1.placeholder === 'Frase Qeqchi';
        inputText1.disabled = false;
        inputText2.disabled = true;
        inputText1.placeholder = isQeqchiPhrase ? 'Frase Español' : 'Frase Qeqchi';
        inputText2.placeholder = isQeqchiPhrase ? 'Frase Qeqchi' : 'Frase Español';
        clearInputs();
    }

    function clearInputs() {
        document.getElementById('inputText1').value = '';
        document.getElementById('inputText2').value = '';
        document.getElementById('suggestedWordsListBox').innerHTML = '';
    }

    function fetchSuggestions() {
        const query = document.getElementById('inputText1').value;

        if (query.length > 0) {
            fetch(`/suggestions?query=${encodeURIComponent(query)}`)
                .then(response => response.json())
                .then(data => {
                    const suggestionsListBox = document.getElementById('suggestedWordsListBox');
                    suggestionsListBox.innerHTML = ''; // Clear previous suggestions
                    data.forEach(suggestion => {
                        const option = document.createElement('option');
                        option.value = suggestion;
                        option.textContent = suggestion;
                        suggestionsListBox.appendChild(option);
                    });
                })
                .catch(error => console.error('Error fetching suggestions:', error));
        } else {
            document.getElementById('suggestedWordsListBox').innerHTML = ''; // Clear suggestions if no text
        }
    }

    function translateText() {
        const inputText1 = document.getElementById('inputText1').value;
        let direction;

        const placeholder = document.getElementById('inputText1').placeholder;
        if (placeholder === 'Español') direction = 'es_to_qeqchi';
        else if (placeholder === 'Qeqchi') direction = 'qeq_to_es';
        else if (placeholder === 'Frase Español') direction = 'fes_to_qeqchi';
        else if (placeholder === 'Frase Qeqchi') direction = 'fqeq_to_es';

        fetch('/translate', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text: inputText1, direction })
        })
            .then(response => response.json())
            .then(data => {
                document.getElementById('inputText2').value = data.translation;
            })
            .catch(error => console.error('Error translating text:', error));
    }

    function fetchWordBank() {
        fetch('/wordbank', {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        })
            .then(response => response.json())
            .then(data => {
                const wordBankWindow = window.open('', '_blank');
                wordBankWindow.document.write('<html><head><title>Banco de Palabras</title></head><body>');
                wordBankWindow.document.write('<h1>Banco de Palabras</h1>');
                wordBankWindow.document.write('<table border="1"><tr><th>Español</th><th>Q\'eqchi\'</th></tr>');

                data.forEach(row => {
                    wordBankWindow.document.write(`<tr><td>${row.Espanol}</td><td>${row.Qeqchi}</td></tr>`);
                });

                wordBankWindow.document.write('</table></body></html>');
                wordBankWindow.document.close();
            })
            .catch(error => console.error('Error fetching word bank:', error));
    }
});
