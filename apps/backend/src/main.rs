use tokio::net::TcpListener;
use the_pound::core;

#[tokio::main]
async fn main() {
    let app_router: axum::Router = core::define_app_routes();
    let listener = TcpListener::bind("127.0.0.1:8000").await.unwrap();
    axum::serve(listener, app_router).await.unwrap();
}
