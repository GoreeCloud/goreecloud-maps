import './search-save-action.css';
import {
  MapsApiError,
  subscribeSearchSaveContext,
  type SearchSaveContext,
} from './api';

let latestContexts: readonly SearchSaveContext[] = [];

const placeLabel = (context: SearchSaveContext): string =>
  context.result.name || context.result.label || 'place';

const markSavedAcrossRenderedCopies = (index: number): void => {
  document.querySelectorAll<HTMLButtonElement>(`[data-search-save-index="${index}"]`).forEach((button) => {
    button.disabled = true;
    button.textContent = 'Saved';
    button.setAttribute('aria-label', `${button.dataset.placeLabel ?? 'Place'} is saved to your private Maps places`);
  });
  document.querySelectorAll<HTMLElement>(`[data-search-save-status="${index}"]`).forEach((status) => {
    status.textContent = 'Saved privately to Maps.';
  });
};

const enhanceSearchResults = (): void => {
  const cards = document.querySelectorAll<HTMLButtonElement>('.result-card[data-result-index]');
  cards.forEach((card) => {
    if (card.dataset.searchSaveEnhanced === 'true') return;
    const index = Number(card.dataset.resultIndex);
    const context = latestContexts[index];
    if (!Number.isInteger(index) || !context) return;

    card.dataset.searchSaveEnhanced = 'true';
    const label = placeLabel(context);
    card.setAttribute('aria-label', `Show ${label} on map`);
    card.title = `Show ${label} on map`;

    const actionHint = document.createElement('span');
    actionHint.className = 'result-card-action-label';
    actionHint.textContent = 'Show on map';
    card.append(actionHint);

    const wrapper = document.createElement('div');
    wrapper.className = 'result-card-group';
    card.before(wrapper);
    wrapper.append(card);

    const actions = document.createElement('div');
    actions.className = 'result-save-actions';

    const saveButton = document.createElement('button');
    saveButton.type = 'button';
    saveButton.className = 'integration-button result-save-button';
    saveButton.textContent = 'Save';
    saveButton.dataset.searchSaveIndex = String(index);
    saveButton.dataset.placeLabel = label;
    saveButton.setAttribute('aria-label', `Save ${label} to your private Maps places`);

    const status = document.createElement('span');
    status.className = 'result-save-status';
    status.dataset.searchSaveStatus = String(index);
    status.setAttribute('role', 'status');
    status.setAttribute('aria-live', 'polite');

    saveButton.addEventListener('click', async () => {
      saveButton.disabled = true;
      saveButton.textContent = 'Saving…';
      status.textContent = 'Saving this result to your private Maps places.';
      try {
        await context.save();
        markSavedAcrossRenderedCopies(index);
      } catch (error) {
        saveButton.disabled = false;
        saveButton.textContent = 'Try again';
        status.textContent = error instanceof MapsApiError
          ? error.message
          : 'Maps could not save this result. No save success is being assumed.';
      }
    });

    actions.append(saveButton, status);
    wrapper.append(actions);
  });
};

subscribeSearchSaveContext((contexts) => {
  latestContexts = contexts;
  queueMicrotask(enhanceSearchResults);
});

const app = document.querySelector<HTMLElement>('#app');
if (app) {
  const observer = new MutationObserver(() => enhanceSearchResults());
  observer.observe(app, { childList: true, subtree: true });
}
