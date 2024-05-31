FROM node:20.5.0-alpine3.18

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

RUN npm install -g serve

CMD ["serve", "-s", "dist"]

EXPOSE 5000
