# Dockerfile for Next.js frontend
FROM node:22.21.1-alpine AS base

# Ensure pnpm is available in the base image
RUN corepack enable && corepack prepare pnpm@latest --activate

# Build stage: install runtime compat, install deps and build the app
FROM base AS builder
RUN apk add --no-cache libc6-compat
WORKDIR /app

# copy lockfile(s) and package manifest first for better layer caching
COPY package.json pnpm-lock.yaml* package-lock.json* ./
RUN pnpm install

# copy source and build
COPY . .
ENV NEXT_TELEMETRY_DISABLED 1
RUN pnpm run build

# Production image, copy only runtime artifacts and run
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs && adduser --system --uid 1001 nextjs

# static/public files
COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown nextjs:nodejs .next

# Automatically leverage output traces to reduce image size
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

CMD ["node", "server.js"]
