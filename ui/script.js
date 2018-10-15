const
    env = document.getElementById('env'),
    add = document.querySelector('#add'),
    restore = document.querySelector('#restore'),
    clear = document.querySelector('#clear'),
    counter = document.querySelector('#counter'),
    current = document.querySelector('#current'),
    file = document.querySelector('#file'),
    save = document.querySelector('#save'),
    saveInput = document.querySelector('#save-input'),
    saveSource = document.querySelector('#save-source'),
    saveForm = document.querySelector('#save-form'),
    search = document.querySelector('#search'),
    template = document.querySelector('#template').innerHTML;

let deleted = [];

// Search
search.addEventListener('keyup', function() {
    const value = this.value.toLowerCase();

    env.querySelectorAll('tr').forEach(function(e) {
        const search = (
            e.querySelector('[name="name"]').value +
            e.querySelector('[name="value"]').value +
            e.querySelector('[name="comment"]').value
        ).toLowerCase();

        if (search.indexOf(value) > -1) {
            e.classList.remove('hidden');
        } else {
            e.classList.add('hidden');
        }
    });
});

// Clear
clear.addEventListener('click', function () {
    clearAll();
    saveSource.value = null;
});

// Load current env
current.addEventListener('click', function() {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', '/env/current', true);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            saveSource.value = null;
            load(JSON.parse(xhr.responseText));
        }
    };
    xhr.send();
});

// Load from file
file.addEventListener('change', function() {
    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/env/file', true);
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            load(JSON.parse(xhr.responseText));
        }
    };

    const reader = new FileReader();
    reader.readAsBinaryString(this.files[0]);
    reader.onload = function() {
        saveSource.value = this.result;
        xhr.send(this.result);
    };

    file.value = '';
}, false);

// Save to file
save.addEventListener('click', function() {
    let isValid = true;

    const elements = env.querySelectorAll('tr');

    if (elements.length === 0) {
        return false;
    }

    const vars = [];
    elements .forEach(function(e) {
        const err = e.querySelector('.error');
        if (err !== null) {
            isValid = false;
            err.focus();
            return false;
        }

        vars.push({
            index: parseInt(e.dataset.index, 10),
            name: e.dataset.name,
            new_name: e.querySelector('[name="name"]').value,
            value: e.querySelector('[name="value"]').value,
            comment: e.querySelector('[name="comment"]').value,
            deleted: e.classList.contains('deleted'),
        });

    });

    if (!isValid) {
        return false;
    }

    saveInput.value = JSON.stringify(vars);
    saveForm.submit();
});

// Delete variable
const remove = function(e){
    e.preventDefault();

    const row = this.parentNode.parentNode;

    deleted.push(row.dataset.index);
    row.classList.add('deleted');

    restore.classList.remove('hidden');
    counter.textContent = deleted.length.toString();
};

// Add new
add.addEventListener('click', function() {
    const v = {
        index: env.querySelectorAll('tr').length + 1,
        name: '',
        value: '',
        comment: '',
    };
    env.innerHTML += renderLine(v);

    const sel = '[data-index="' + v.index + '"]';
    env.querySelector(sel).querySelector('[name="name"]').focus();
});

// Restore deleted variables
restore.addEventListener('click', function() {
    if (deleted.length === 0) {
        return
    }

    const index = deleted.pop();
    const s = '[data-index="' + index +'"]';
    env.querySelector(s).classList.remove('deleted');

    if (deleted.length === 0) {
        restore.classList.add('hidden');
    }
    counter.textContent = deleted.length.toString();
});

// Validation
const validate = function(e) {
    e.classList.toggle('error', !e.validity.valid);
};

// Watch for list changes
watchForChanges = function() {
    const elements = env.querySelectorAll('.remove a');

    elements.forEach(function(e) {
        e.removeEventListener('click', remove);
        e.addEventListener('click', remove);
    });

    env.querySelectorAll('.validate').forEach(function(e) {
        e.removeEventListener('keyup', function() {
            validate(e)
        });
        e.addEventListener('keyup', function() {
            validate(e)
        });
    })
};

const envChangeCallback = function(mutationsList, observer) {
    for (const mutation of mutationsList) {
        if (mutation.type === 'childList') {
            watchForChanges();
        }
    }
};

const observer = new MutationObserver(envChangeCallback);
observer.observe(env, {childList: true});

// Render list
function load(variables) {
    clearAll();

    variables.forEach(function(v, index) {
        if (v.index === undefined) {
            v.index = index + 1;
        }

        env.innerHTML += renderLine(v);
    });
}

function renderLine(variable) {
    let line = template;

    for (const prop in variable) {
        const placeholder = new RegExp('{{' + prop + '}}', 'g');
        if (variable.hasOwnProperty(prop)) {
            line = line.replace(placeholder, variable[prop]);
        }
    }

    return line;
}

function clearAll() {
    env.innerHTML = '';
    restore.classList.add('hidden');
    counter.textContent = '';
    search.value = '';
    deleted = [];
}
