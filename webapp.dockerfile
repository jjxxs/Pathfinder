FROM node:8 as builder
COPY webui/public /webui/public
COPY webui/src /webui/src
COPY webui/package.json /webui/package.json
COPY webui/tsconfig.json /webui/tsconfig.json
WORKDIR /webui
RUN npm install && npm run build

FROM nginx:stable
COPY --from=builder /webui/build /usr/share/nginx/html
