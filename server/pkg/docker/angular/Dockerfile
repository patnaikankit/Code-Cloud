FROM node:20.5.0-alpine3.18

WORKDIR /app

COPY package*.json ./

RUN npm install -g @angular/cli

RUN npm install

COPY . .

RUN ng build --prod

FROM nginx:alpine

COPY --from=0 /app/dist/[your-angular-app-name] /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
