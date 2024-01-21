# Stage 1: Build the application
FROM node:lts-alpine as build-stage
WORKDIR /app
COPY view/package*.json ./
RUN npm install
COPY view/ .
RUN npm run build

# Stage 2: Serve the application from Nginx
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]