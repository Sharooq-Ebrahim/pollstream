const API_BASE = window.location.origin;
const WS_BASE = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
const HOST = window.location.host;

let socket = null;
let currentPollId = null;
let currentPollData = null;


const pollSelector = document.getElementById('poll-selector');
const livePoll = document.getElementById('live-poll');
const pollsList = document.getElementById('polls-list');
const pollQuestion = document.getElementById('poll-question');
const pollOptions = document.getElementById('poll-options');
const totalVotesEl = document.getElementById('total-votes');
const backBtn = document.getElementById('back-btn');


document.addEventListener('DOMContentLoaded', fetchPolls);

async function fetchPolls() {
    try {
        const response = await fetch(`${API_BASE}/poll/getAll`);
        const polls = await response.json();
        
        pollsList.innerHTML = '';
        if (polls.length === 0) {
            pollsList.innerHTML = '<p class="text-dim">No active polls found.</p>';
            return;
        }

        polls.forEach(poll => {
            const optionsCount = poll.options ? poll.options.length : 0;
            console.log('options',poll.options);
            const div = document.createElement('div');
            div.className = 'poll-item';
            div.innerHTML = `
                <div class="poll-info">
                    <h3>${poll.question}</h3>
                    <p>${optionsCount} options</p>
                </div>
                <span class="view-btn">View &rarr;</span>
            `;
            div.onclick = () => joinPoll(poll.id);
            pollsList.appendChild(div);
        });
    } catch (err) {
        console.error('Error fetching polls:', err);
        pollsList.innerHTML = '<p class="text-dim">Error loading polls. Is the server running?</p>';
    }
}

function joinPoll(id) {
    currentPollId = id;
    pollSelector.classList.add('hidden');
    livePoll.classList.remove('hidden');
    
    connectWebSocket(id);
}

function connectWebSocket(id) {
    if (socket) socket.close();

    const url = `${WS_BASE}${HOST}/ws/poll?id=${id}`;
    socket = new WebSocket(url);

    socket.onmessage = (event) => {
        const poll = JSON.parse(event.data);
        if (poll.id === currentPollId) {
            updateUI(poll);
        }
    };

    socket.onclose = () => {
        console.log('WebSocket disconnected');
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

function updateUI(poll) {
    currentPollData = poll;
    pollQuestion.innerText = poll.question;
    pollOptions.innerHTML = '';
    
   
    const votedOptionId = localStorage.getItem(`voted_${poll.id}`);
    const isVoted = !!votedOptionId;

    const options = (poll.options || []).sort((a, b) => a.id.localeCompare(b.id));
    let totalVotes = 0;
    options.forEach(opt => totalVotes += opt.votes);
    
    totalVotesEl.innerText = isVoted ? `Total Votes: ${totalVotes}` : 'Cast your vote to see results';

    options.forEach(option => {
        const percentage = totalVotes > 0 ? (option.votes / totalVotes) * 100 : 0;
        const votedForThis = option.id === votedOptionId;

        const wrapper = document.createElement('div');
        wrapper.className = 'option-wrapper';
        
        wrapper.innerHTML = `
            <div class="progress-bar ${isVoted ? '' : 'hidden-data'}" style="width: ${isVoted ? percentage : 0}%"></div>
            <button class="option-btn ${votedForThis ? 'voted' : ''}" 
                    ${isVoted ? 'disabled' : ''} 
                    onclick="vote('${poll.id}', '${option.id}')">
                <div class="option-data">
                    <span class="option-text">
                        ${option.text}
                        ${votedForThis ? '<span class="voted-marker">✓ YOUR VOTE</span>' : ''}
                    </span>
                    <span class="vote-count ${isVoted ? '' : 'hidden-data'}">
                        ${isVoted ? option.votes : '?'}
                    </span>
                </div>
            </button>
        `;
        pollOptions.appendChild(wrapper);
    });
}

async function vote(pollId, optionId) {
    if (localStorage.getItem(`voted_${pollId}`)) {
        showToast('You have already voted in this poll', true);
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/poll/vote`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'poll_id': pollId,
                'option_id': optionId
            })
        });

        if (response.ok) {
  
            localStorage.setItem(`voted_${pollId}`, optionId);
            showToast('Vote counted!');
           
            if (currentPollData && currentPollData.id === pollId) {
                updateUI(currentPollData);
            }
        } else {
            showToast('Error casting vote', true);
        }
    } catch (err) {
        console.error('Vote error:', err);
    }
}

function showToast(message, isError = false) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'toast';
    if (isError) toast.style.backgroundColor = 'var(--danger)';
    toast.innerText = message;
    
    container.appendChild(toast);
    setTimeout(() => toast.remove(), 3000);
}

const showCreateBtn = document.getElementById('show-create-poll');
const closeCreateBtn = document.getElementById('close-create-btn');
const createPollSection = document.getElementById('create-poll-section');
const addOptionBtn = document.getElementById('add-option-btn');
const optionsInputs = document.getElementById('options-inputs').querySelector('.option-input-wrapper');
const submitPollBtn = document.getElementById('submit-poll-btn');

showCreateBtn.onclick = () => {
    pollSelector.classList.add('hidden');
    containerToggle(createPollSection, true);
};

closeCreateBtn.onclick = () => {
    containerToggle(createPollSection, false);
    pollSelector.classList.remove('hidden');
};

addOptionBtn.onclick = () => {
    const input = document.createElement('input');
    input.type = 'text';
    input.className = 'option-input';
    input.placeholder = `Option ${optionsInputs.children.length + 1}`;
    optionsInputs.appendChild(input);
};

submitPollBtn.onclick = async () => {
    const question = document.getElementById('poll-question-input').value;
    const optionInputs = document.querySelectorAll('.option-input');
    const options = Array.from(optionInputs)
        .map(input => input.value.trim())
        .filter(val => val !== '')
        .map((text, index) => ({
            id: `opt-${Date.now()}-${index}`,
            text: text,
            votes: 0
        }));

    if (!question || options.length < 2) {
        showToast('Please provide a question and at least 2 options', true);
        return;
    }

    const newPoll = {
        id: `poll-${Date.now()}`,
        question: question,
        options: options
    };

    try {
        const response = await fetch(`${API_BASE}/poll/create`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newPoll)
        });

        if (response.ok) {
            showToast('Poll created successfully!');
            containerToggle(createPollSection, false);
            pollSelector.classList.remove('hidden');
            fetchPolls();
            document.getElementById('poll-question-input').value = '';
            optionsInputs.innerHTML = `
                <input type="text" class="option-input" placeholder="Option 1">
                <input type="text" class="option-input" placeholder="Option 2">
            `;
        } else {
            showToast('Error creating poll', true);
        }
    } catch (err) {
        console.error('Error creating poll:', err);
    }
};

function containerToggle(el, show) {
    if (show) el.classList.remove('hidden');
    else el.classList.add('hidden');
}

backBtn.onclick = () => {
    if (socket) socket.close();
    currentPollId = null;
    currentPollData = null;
    livePoll.classList.add('hidden');
    pollSelector.classList.remove('hidden');
    fetchPolls();
};
