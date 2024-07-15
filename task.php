<?php


//php -S 127.0.0.1:8088 -t .
if (isset($_GET["action"])&&($_GET["action"] === "queryResult")) {
    $d = json_encode($_POST);
    file_put_contents(__DIR__."/queryResult.json", $d);
    echo "hahahahha";
    die;
}


$data = [

	"processingImage" => [
		["metaId" => strval(rand(1000,999999)), "imageUrl" => "https://cdn.fyh88.com/upload/1-ee580bff5d2fbf4fb28329a8ea8daedf.jpeg" ],
		["metaId" => strval(rand(1000,999999)), "imageUrl" => "https://cdn.fyh88.com/upload/1-648964d3f0fd9cc253c28acb110fc52e.jpeg" ],
	
	],
	"queryImage" => [
	        ["queryId" => strval(rand(1000,999999)),"imageUrl" => "https://cdn.fyh88.com/upload/1-25c79166eddfd4161f9d8a7f56be4971.jpeg"],
	        ["queryId" => strval(rand(1000,999999)),"imageUrl" => "https://cdn.fyh88.com/upload/1-36d06ce80f325e1ea8c76b50720a0c36.jpeg"],
	]
];


header('Content-Type: application/json');
echo json_encode($data);


