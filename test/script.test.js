import { beforeAll, afterAll, describe, expect, it, vi } from 'vitest';

// The site script (web-assets/js/script.js) is a plain browser script that runs
// side effects on load (registering timers and DOM listeners) and exposes a set
// of functions on `window.GoLucky`. We import it once inside a jsdom environment
// with fake timers so the periodic intervals it schedules never actually fire.
let GoLucky;

beforeAll(async () => {
    vi.useFakeTimers();
    await import('../web-assets/js/script.js');
    GoLucky = window.GoLucky;
});

afterAll(() => {
    vi.useRealTimers();
});

describe('module surface', () => {
    it('exposes the expected functions on window.GoLucky', () => {
        expect(GoLucky).toBeTypeOf('object');
        for (const fn of [
            'generateRandomCosmicNumbers',
            'generateRandomLuckyBall',
            'debounce',
            'getDomainFromUrl',
            'getButtonType',
            'getSourceSection',
        ]) {
            expect(GoLucky[fn], `GoLucky.${fn}`).toBeTypeOf('function');
        }
    });
});

describe('generateRandomCosmicNumbers', () => {
    it('returns 5 unique numbers in [1, 48], sorted ascending', () => {
        // Run many times to exercise the randomness rather than a single sample.
        for (let i = 0; i < 500; i++) {
            const nums = GoLucky.generateRandomCosmicNumbers();

            expect(nums).toHaveLength(5);
            expect(new Set(nums).size).toBe(5); // all unique

            for (const n of nums) {
                expect(Number.isInteger(n)).toBe(true);
                expect(n).toBeGreaterThanOrEqual(1);
                expect(n).toBeLessThanOrEqual(48);
            }

            const sorted = [...nums].sort((a, b) => a - b);
            expect(nums).toEqual(sorted);
        }
    });
});

describe('generateRandomLuckyBall', () => {
    it('returns an integer in [1, 18]', () => {
        for (let i = 0; i < 500; i++) {
            const ball = GoLucky.generateRandomLuckyBall();
            expect(Number.isInteger(ball)).toBe(true);
            expect(ball).toBeGreaterThanOrEqual(1);
            expect(ball).toBeLessThanOrEqual(18);
        }
    });
});

describe('getDomainFromUrl', () => {
    it('extracts the hostname from a valid URL', () => {
        expect(GoLucky.getDomainFromUrl('https://github.com/mrz1836/go-lucky')).toBe('github.com');
        expect(GoLucky.getDomainFromUrl('http://example.com:8080/path?q=1')).toBe('example.com');
    });

    it('returns "unknown" for an unparseable URL', () => {
        expect(GoLucky.getDomainFromUrl('not a url')).toBe('unknown');
        expect(GoLucky.getDomainFromUrl('')).toBe('unknown');
    });
});

describe('getButtonType', () => {
    const makeEl = (tag, className = '') => {
        const el = document.createElement(tag);
        if (className) el.className = className;
        return el;
    };

    it('classifies primary and secondary buttons by class', () => {
        expect(GoLucky.getButtonType(makeEl('a', 'btn btn-primary'))).toBe('primary');
        expect(GoLucky.getButtonType(makeEl('a', 'btn btn-secondary'))).toBe('secondary');
    });

    it('classifies strategy tabs', () => {
        expect(GoLucky.getButtonType(makeEl('div', 'strategy-tab'))).toBe('strategy-tab');
    });

    it('falls back to tag name for plain button and anchor elements', () => {
        expect(GoLucky.getButtonType(makeEl('button'))).toBe('button');
        expect(GoLucky.getButtonType(makeEl('a'))).toBe('link');
    });

    it('returns "unknown" for unrecognized elements', () => {
        expect(GoLucky.getButtonType(makeEl('span'))).toBe('unknown');
    });
});

describe('getSourceSection', () => {
    it('returns the known section class when present', () => {
        const section = document.createElement('section');
        section.className = 'features some-other-class';
        const child = document.createElement('button');
        section.appendChild(child);

        expect(GoLucky.getSourceSection(child)).toBe('features');
    });

    it('falls back to the section id when no known class matches', () => {
        const section = document.createElement('section');
        section.id = 'custom-section';
        const child = document.createElement('a');
        section.appendChild(child);

        expect(GoLucky.getSourceSection(child)).toBe('custom-section');
    });

    it('falls back to the tag name for a bare container', () => {
        const footer = document.createElement('footer');
        const child = document.createElement('a');
        footer.appendChild(child);

        // "footer" is in the known-class list, but here it is a tag with no class,
        // so classification falls through to the lowercased tag name.
        expect(GoLucky.getSourceSection(child)).toBe('footer');
    });

    it('returns "unknown" when the element is not inside any container', () => {
        const orphan = document.createElement('a');
        expect(GoLucky.getSourceSection(orphan)).toBe('unknown');
    });
});

describe('debounce', () => {
    it('invokes the wrapped function once after the wait window', () => {
        const spy = vi.fn();
        const debounced = GoLucky.debounce(spy, 100);

        debounced();
        debounced();
        debounced();

        expect(spy).not.toHaveBeenCalled();

        vi.advanceTimersByTime(100);
        expect(spy).toHaveBeenCalledTimes(1);
    });

    it('passes the most recent arguments through to the wrapped function', () => {
        const spy = vi.fn();
        const debounced = GoLucky.debounce(spy, 50);

        debounced('first');
        debounced('second');

        vi.advanceTimersByTime(50);
        expect(spy).toHaveBeenCalledTimes(1);
        expect(spy).toHaveBeenCalledWith('second');
    });
});
