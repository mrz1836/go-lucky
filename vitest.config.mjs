import { defineConfig } from 'vitest/config';

export default defineConfig({
    test: {
        // The site script manipulates the DOM directly, so run tests in jsdom.
        environment: 'jsdom',
        include: ['test/**/*.test.{js,mjs}'],
        coverage: {
            provider: 'v8',
            reporter: ['text', 'html'],
            include: ['web-assets/js/**/*.js'],
        },
    },
});
