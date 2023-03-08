+++
date = "2023-02-02-10T06:43:48+02:00"
title = "PHP Client for the Greypot Studio API"
draft = false
weight = 50
description = "PHP Client for the Greypot Studio API"
toc = true
bref = "Client.PHP"
+++


## Laravel

You can create a class in `app\Http\Client` with the following content

```php
<?php 

namespace App\Http\Client;

use Illuminate\Support\Facades\Http;

class GreypotClient 
{
    private $greypotBaseURL = "https://greypot-studio.fly.dev/_studio";

    private $reportStoreDirectory = null;

    public function __construct($reportStoreDirectory = 'data')
    {
        $this->reportStoreDirectory = $reportStoreDirectory;
    }

    public function generatePDF($template, $data) 
    {
        $response = Http::withUrlParameters([
            'endpoint' => $this->greypotBaseURL,
            'reportId' => uniqid(),
        ])->post('{+endpoint}/generate/pdf/{reportId}', [
            'Name' => 'reportId',
            'Template' => $template,
            'Data' => $data,
        ], [
            'Content-Type' => 'application/json'
        ]);

        if ($response->ok()) {
            $fileDataBase64 = $response->json('data');
            return base64_decode($fileDataBase64);
        }
        throw new \Exception(sprintf('failed to generate report from template=%s got response=%s', $template, $response->body() ));
    }

    public function generateBulkPDF($template, $entries = [])
    {
        $response = Http::withUrlParameters([
            'endpoint' => $this->greypotBaseURL,
            'reportId' => uniqid(),
        ])->post('{+endpoint}/generate/bulk/pdf/{reportId}', [
            'Name' => 'reportId',
            'Template' => $template,
            'Data' => $entries,
        ], [
            'Content-Type' => 'application/json'
        ]);

        if ($response->ok()) {
            $reports = $response->json('reports');
            return $reports;
        }
        throw new \Exception(sprintf('failed to generate reports from template=%s got response=%s', $template, $response->body() ));
    }
}
```

Example Usage

```php
<?php
namespace App\Http\Controllers;
use Illuminate\Foundation\Auth\Access\AuthorizesRequests;
use Illuminate\Foundation\Validation\ValidatesRequests;
use Illuminate\Routing\Controller as BaseController;

use App\Http\Client\GreypotClient;

class ReportController extends BaseController
{

    private $greypot;
    
    public function __construct(GreypotClient $greypot)
    {
        $this->greypot = $greypot;
    }

    public function download()
    {
        $data = [
            'contributors' => [
                [
                    'name' => 'XXXXX',
                    'phone' => 'XXXXX',
                    'email' => 'XXXXX',
                    'amount' => 'XXXXX',
                    'date' => 'XXXXX',
                ]
            ]
        ];

        $template = file_get_contents(resource_path('report-template/contributors.html'));
        
        $reportData = $this->greypot->generatePDF($template, $data);

        return response()->streamDownload(function() use ($reportData) {
                echo $reportData;
            }, 'report.pdf');
    }    
}
```


## Use Guzzle HTTP

## 